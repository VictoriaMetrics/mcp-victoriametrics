package resources

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/config"
	"github.com/VictoriaMetrics/mcp-victoriametrics/cmd/mcp-victoriametrics/utils"
)

//go:embed vm/content vmsite/content/blog
var DocsDir embed.FS

const (
	docsURIPrefix              = "docs://"
	maxMarkdownDescriptionSize = 4096
)

var (
	searchIndex bleve.Index
	resources   map[string]mcp.Resource
	contents    map[string]mcp.ResourceContents
)

func RegisterDocsResources(s *server.MCPServer, cfg *config.Config) {
	var err error
	mapping := bleve.NewIndexMapping()
	if searchIndex, err = bleve.NewMemOnly(mapping); err != nil {
		log.Fatal(fmt.Errorf("error creating index: %w", err))
	}

	docFiles, err := ListDocFiles()
	if err != nil {
		log.Fatal(fmt.Errorf("error listing docs files: %w", err))
	}
	resources = make(map[string]mcp.Resource, len(docFiles))
	contents = make(map[string]mcp.ResourceContents, len(docFiles))
	for _, docFile := range docFiles {
		resourceURI := fmt.Sprintf("%s%s#%d", docsURIPrefix, docFile.Path, docFile.ChunkNum)
		resource := mcp.NewResource(
			resourceURI,
			docFile.Name,
			mcp.WithMIMEType("text/markdown"),
			mcp.WithResourceDescription(docFile.Content[:min(len(docFile.Content), maxMarkdownDescriptionSize)]),
		)
		if !cfg.IsResourcesDisabled() {
			s.AddResource(resource, docResourcesHandler)
		}
		resources[resourceURI] = resource
		contents[resourceURI] = mcp.TextResourceContents{
			URI:      resourceURI,
			MIMEType: "text/markdown",
			Text:     docFile.Content,
		}
		if err = searchIndex.Index(resourceURI, docFile); err != nil {
			log.Fatal(fmt.Errorf("error indexing file %s: %w", docFile.Path, err))
		}
	}
}

func SearchDocResources(query string, limit int) ([]mcp.Resource, error) {
	searchQuery := bleve.NewMatchQuery(query)
	searchQuery.Fuzziness = 1
	searchRequest := bleve.NewSearchRequest(searchQuery)
	searchRequest.Size = limit
	searchResults, err := searchIndex.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error searching index: %w", err)
	}
	if searchResults.Total == 0 {
		return nil, fmt.Errorf("no results found for query: %s", query)
	}
	results := make([]mcp.Resource, 0)
	for _, hit := range searchResults.Hits[:limit] {
		resource, ok := resources[hit.ID]
		if !ok {
			continue
		}
		results = append(results, resource)
	}
	return results, nil
}

func docResourcesHandler(_ context.Context, rrr mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	content, err := GetDocResourceContent(rrr.Params.URI)
	if err != nil {
		return nil, fmt.Errorf("error getting doc resource content: %w", err)
	}
	return []mcp.ResourceContents{content}, nil
}

func GetDocResourceContent(uri string) (mcp.ResourceContents, error) {
	content, ok := contents[uri]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", uri)
	}
	return content, nil
}

func GetDocFileContent(path string) (string, error) {
	file, err := fs.ReadFile(DocsDir, path)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %w", path, err)
	}
	return string(file), nil
}

type DocFileInfo struct {
	Path     string `json:"path"`
	ChunkNum int    `json:"chunk_num"`
	Content  string `json:"content"`
	Name     string `json:"name"`
}

func ListDocFiles() ([]DocFileInfo, error) {
	docs := make([]DocFileInfo, 0)
	for _, rootDir := range []string{"vm", "vmsite"} {
		docFiles, err := utils.Glob(DocsDir, rootDir, func(s string) bool {
			return strings.ToLower(filepath.Ext(s)) == ".md"
		})
		if err != nil {
			return nil, fmt.Errorf("error reading docs directory: %w", err)
		}
		for _, path := range docFiles {
			if !strings.HasSuffix(strings.ToLower(path), ".md") {
				continue
			}
			content, err := GetDocFileContent(path)
			if err != nil {
				return nil, fmt.Errorf("error reading file %s: %w", path, err)
			}

			chunks, err := splitMarkdown(content)
			if err != nil {
				return nil, fmt.Errorf("error splitting file %s: %w", path, err)
			}

			for chunkNum, chunkContent := range chunks {
				name := ""
				for line := range strings.Lines(chunkContent) {
					if strings.TrimSpace(line) == "" {
						continue
					}
					if !strings.HasPrefix(line, "#") {
						break
					}
					title := strings.TrimSpace(strings.Trim(line, "# "))
					name = fmt.Sprintf("%s / %s", name, title)
				}
				name = strings.Trim(name, "/ ")
				if name == "" {
					name = path
				}

				docs = append(docs, DocFileInfo{
					Path:     path,
					ChunkNum: chunkNum,
					Content:  chunkContent,
					Name:     name,
				})
			}
		}
	}
	return docs, nil
}
