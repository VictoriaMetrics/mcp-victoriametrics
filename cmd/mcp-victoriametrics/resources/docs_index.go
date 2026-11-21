package resources

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
)

//go:embed all:docsindex
var docsIndexDir embed.FS

const (
	embeddedDocsIndexArchive = "docsindex/docs-index.tar.gz"
	docsIndexDirPrefix       = "mcp-victoriametrics-docs-index-"
	docsHashKey              = "docs_hash"
	docsIndexBatchSize       = 64
	docsIndexFormat          = "v1"
)

func BuildDocsIndex(archivePath string) error {
	docFiles, err := ListDocFiles()
	if err != nil {
		return fmt.Errorf("error listing docs files: %w", err)
	}
	tmpDir, err := os.MkdirTemp("", docsIndexDirPrefix+"build-")
	if err != nil {
		return fmt.Errorf("error creating temporary directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	indexDir := filepath.Join(tmpDir, "index")
	if err := buildDocsIndex(indexDir, docFiles); err != nil {
		return err
	}
	return writeDocsIndexArchive(indexDir, archivePath)
}

func buildDocsIndex(path string, docFiles []DocFileInfo) error {
	index, err := bleve.New(path, bleve.NewIndexMapping())
	if err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	if err := indexDocs(index, docFiles); err != nil {
		_ = index.Close()
		return err
	}
	if err := index.SetInternal([]byte(docsHashKey), []byte(docsHash(docFiles))); err != nil {
		_ = index.Close()
		return fmt.Errorf("error storing docs hash: %w", err)
	}
	return index.Close()
}

func openDocsIndex(docFiles []DocFileInfo) bleve.Index {
	hash := docsHash(docFiles)
	archive, err := docsIndexDir.ReadFile(embeddedDocsIndexArchive)
	if err == nil {
		dir := filepath.Join(os.TempDir(), docsIndexDirPrefix+hash[:16])
		index, err := openArchivedDocsIndex(archive, dir, hash)
		if err == nil {
			return index
		}
		slog.Warn("cannot use embedded docs index, building it in memory", "error", err)
	}
	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		log.Fatal(fmt.Errorf("error creating index: %w", err))
	}
	if err := indexDocs(index, docFiles); err != nil {
		log.Fatal(err)
	}
	return index
}

func openArchivedDocsIndex(archive []byte, dir, hash string) (bleve.Index, error) {
	if index, err := openPrebuiltDocsIndex(dir, hash); err == nil {
		return index, nil
	}
	if err := os.RemoveAll(dir); err != nil {
		return nil, fmt.Errorf("error removing outdated index: %w", err)
	}
	tmpDir, err := os.MkdirTemp(filepath.Dir(dir), filepath.Base(dir)+".tmp-")
	if err != nil {
		return nil, fmt.Errorf("error creating index directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()
	if err := extractDocsIndexArchive(archive, tmpDir); err != nil {
		return nil, err
	}
	if err := os.Rename(tmpDir, dir); err != nil {
		if _, statErr := os.Stat(dir); statErr != nil {
			return nil, fmt.Errorf("error moving index: %w", err)
		}
	}
	return openPrebuiltDocsIndex(dir, hash)
}

func openPrebuiltDocsIndex(path, hash string) (bleve.Index, error) {
	index, err := bleve.OpenUsing(path, map[string]any{"read_only": true})
	if err != nil {
		return nil, err
	}
	stored, err := index.GetInternal([]byte(docsHashKey))
	if err != nil {
		_ = index.Close()
		return nil, fmt.Errorf("error reading docs hash: %w", err)
	}
	if string(stored) != hash {
		_ = index.Close()
		return nil, fmt.Errorf("index was built for different docs")
	}
	return index, nil
}

func writeDocsIndexArchive(indexDir, archivePath string) (err error) {
	tmpPath := archivePath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("error creating archive: %w", err)
	}
	defer func() {
		if err != nil {
			_ = f.Close()
			_ = os.Remove(tmpPath)
		}
	}()
	gz, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	tw := tar.NewWriter(gz)
	err = filepath.WalkDir(indexDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || path == indexDir {
			return err
		}
		rel, err := filepath.Rel(indexDir, path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(rel)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = src.Close() }()
		_, err = io.Copy(tw, src)
		return err
	})
	if err != nil {
		return fmt.Errorf("error writing archive: %w", err)
	}
	if err := errors.Join(tw.Close(), gz.Close(), f.Close()); err != nil {
		return fmt.Errorf("error writing archive: %w", err)
	}
	return os.Rename(tmpPath, archivePath)
}

func extractDocsIndexArchive(archive []byte, dir string) error {
	gz, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return fmt.Errorf("error reading archive: %w", err)
	}
	tr := tar.NewReader(gz)
	for {
		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error reading archive: %w", err)
		}
		if !filepath.IsLocal(header.Name) {
			return fmt.Errorf("unexpected path in archive: %s", header.Name)
		}
		target := filepath.Join(dir, filepath.FromSlash(header.Name))
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := extractFile(tr, target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected entry type in archive: %s", header.Name)
		}
	}
}

func extractFile(r io.Reader, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

func indexDocs(index bleve.Index, docFiles []DocFileInfo) error {
	for len(docFiles) > 0 {
		n := min(len(docFiles), docsIndexBatchSize)
		batch := index.NewBatch()
		for _, docFile := range docFiles[:n] {
			if err := batch.Index(docResourceURI(docFile), docFile); err != nil {
				return fmt.Errorf("error indexing file %s: %w", docFile.Path, err)
			}
		}
		if err := index.Batch(batch); err != nil {
			return fmt.Errorf("error indexing docs: %w", err)
		}
		docFiles = docFiles[n:]
	}
	return nil
}

func docsHash(docFiles []DocFileInfo) string {
	h := sha256.New()
	_, _ = io.WriteString(h, docsIndexFormat)
	for _, docFile := range docFiles {
		for _, v := range []string{docResourceURI(docFile), docFile.Name, docFile.Content} {
			_, _ = io.WriteString(h, v)
			_, _ = h.Write([]byte{0})
		}
	}
	return hex.EncodeToString(h.Sum(nil))
}
