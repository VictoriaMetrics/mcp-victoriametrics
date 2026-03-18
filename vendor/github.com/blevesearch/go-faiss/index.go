package faiss

/*
#include <stdlib.h>
#include <faiss/c_api/Index_c.h>
#include <faiss/c_api/IndexIVF_c.h>
#include <faiss/c_api/IndexIVF_c_ex.h>
#include <faiss/c_api/Index_c_ex.h>
#include <faiss/c_api/impl/AuxIndexStructures_c.h>
#include <faiss/c_api/index_factory_c.h>
#include <faiss/c_api/MetaIndexes_c.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"sort"
	"unsafe"
)

// Index is a Faiss index.
//
// Note that some index implementations do not support all methods.
// Check the Faiss wiki to see what operations an index supports.
type Index interface {
	// D returns the dimension of the indexed vectors.
	D() int

	// IsTrained returns true if the index has been trained or does not require
	// training.
	IsTrained() bool

	// Ntotal returns the number of indexed vectors.
	Ntotal() int64

	// MetricType returns the metric type of the index.
	MetricType() int

	// Train trains the index on a representative set of vectors.
	Train(x []float32) error

	// Add adds vectors to the index.
	Add(x []float32) error

	// AddWithIDs is like Add, but stores xids instead of sequential IDs.
	AddWithIDs(x []float32, xids []int64) error

	// Returns true if the index is an IVF index.
	IsIVFIndex() bool

	// Applicable only to IVF indexes: Returns a slice where each index represents
	// a cluster (list) ID and the value is the count of selected vectors belonging
	// to that cluster. Only vectors specified by the given Selector are considered.
	ObtainClusterVectorCountsFromIVFIndex(include Selector, nlist int) ([]int64, error)

	// Applicable only to IVF indexes: Returns the centroid IDs in the selector in
	// decreasing order of proximity to query 'x' and their distance from 'x'
	ObtainClustersWithDistancesFromIVFIndex(x []float32, centroids Selector, numCentroids int64) (
		[]int64, []float32, error)

	// Applicable only to IVF indexes: Returns the top k centroid cardinalities and
	// their vectors in chosen order (descending or ascending)
	ObtainKCentroidCardinalitiesFromIVFIndex(limit int, descending bool) ([]uint64, [][]float32, error)

	// Search queries the index with the vectors in x.
	// Returns the IDs of the k nearest neighbors for each query vector and the
	// corresponding distances.
	Search(x []float32, k int64) (distances []float32, labels []int64, err error)

	SearchWithoutIDs(x []float32, k int64, exclude Selector, params json.RawMessage) (distances []float32,
		labels []int64, err error)

	SearchWithIDs(x []float32, k int64, include Selector, params json.RawMessage) (distances []float32,
		labels []int64, err error)

	// Applicable only to IVF indexes: Search clusters whose IDs are in eligibleCentroidIDs
	SearchClustersFromIVFIndex(eligibleCentroidIDs []int64, centroidDis []float32, centroidsToProbe int,
		x []float32, k int64, include Selector, params json.RawMessage) ([]float32, []int64, error)

	Reconstruct(key int64) ([]float32, error)

	ReconstructBatch(keys []int64, recons []float32) ([]float32, error)

	MergeFrom(other Index, add_id int64) error

	// RangeSearch queries the index with the vectors in x.
	// Returns all vectors with distance < radius.
	RangeSearch(x []float32, radius float32) (*RangeSearchResult, error)

	// DistCompute computes the distance between the query vector and the vectors specified by ids.
	DistCompute(x []float32, labels []int64) ([]float32, error)

	// Reset removes all vectors from the index.
	Reset() error

	// RemoveIDs removes the vectors specified by sel from the index.
	// Returns the number of elements removed and error.
	RemoveIDs(sel *IDSelector) (int, error)

	// Close frees the memory used by the index.
	Close()

	// consults the C++ side to get the size of the index
	Size() uint64

	cPtr() *C.FaissIndex
}

type faissIndex struct {
	idx *C.FaissIndex
}

func (idx *faissIndex) cPtr() *C.FaissIndex {
	return idx.idx
}

func (idx *faissIndex) Size() uint64 {
	size := C.faiss_Index_size(idx.idx)
	return uint64(size)
}

func (idx *faissIndex) D() int {
	return int(C.faiss_Index_d(idx.idx))
}

func (idx *faissIndex) IsTrained() bool {
	return C.faiss_Index_is_trained(idx.idx) != 0
}

func (idx *faissIndex) Ntotal() int64 {
	return int64(C.faiss_Index_ntotal(idx.idx))
}

func (idx *faissIndex) MetricType() int {
	return int(C.faiss_Index_metric_type(idx.idx))
}

func (idx *faissIndex) Train(x []float32) error {
	n := len(x) / idx.D()
	if c := C.faiss_Index_train(idx.idx, C.idx_t(n), (*C.float)(&x[0])); c != 0 {
		return getLastError()
	}
	return nil
}

func (idx *faissIndex) Add(x []float32) error {
	n := len(x) / idx.D()
	if c := C.faiss_Index_add(idx.idx, C.idx_t(n), (*C.float)(&x[0])); c != 0 {
		return getLastError()
	}
	return nil
}

func (idx *faissIndex) ObtainClusterVectorCountsFromIVFIndex(includedVectors Selector, nlist int) ([]int64, error) {
	// Applicable only to IVF indexes
	ivfPtr := C.faiss_IndexIVF_cast(idx.cPtr())
	if ivfPtr == nil {
		return nil, fmt.Errorf("index is not of ivf type")
	}
	// Creating a slice to hold the count of vectors per cluster
	// Since we have nlist clusters, we create a slice of size nlist
	// listCount[i] will hold the count of vectors in cluster i
	listCount := make([]int64, nlist)
	// Creating a FAISS selector based on the include bitmap.
	params, err := NewStandardSearchParams(includedVectors)
	if err != nil {
		return nil, err
	}
	defer params.Delete()
	// Calling the C function to populate listCount
	// with the count of vectors per cluster, considering only
	// the vectors specified in the include selector.
	if c := C.faiss_IndexIVF_list_vector_count(
		ivfPtr,
		(*C.idx_t)(unsafe.Pointer(&listCount[0])),
		C.size_t(nlist),
		params.sp,
	); c != 0 {
		return nil, getLastError()
	}
	return listCount, nil
}

func (idx *faissIndex) IsIVFIndex() bool {
	if ivfIdx := C.faiss_IndexIVF_cast(idx.cPtr()); ivfIdx == nil {
		return false
	}
	return true
}

func (idx *faissIndex) ObtainClustersWithDistancesFromIVFIndex(x []float32, includedCentroids Selector, numCentroids int64) (
	[]int64, []float32, error) {
	// Applicable only to IVF indexes
	ivfPtr := C.faiss_IndexIVF_cast(idx.cPtr())
	if ivfPtr == nil {
		return nil, nil, fmt.Errorf("index is not of ivf type")
	}
	params, err := NewStandardSearchParams(includedCentroids)
	if err != nil {
		return nil, nil, err
	}
	defer params.Delete()

	// Populate these with the centroids and their distances.
	centroids := make([]int64, numCentroids)
	centroidDistances := make([]float32, numCentroids)

	n := len(x) / idx.D()

	if c := C.faiss_IndexIVF_search_closest_eligible_centroids(
		ivfPtr,
		(C.idx_t)(n),
		(*C.float)(&x[0]),
		(C.idx_t)(numCentroids),
		(*C.float)(&centroidDistances[0]),
		(*C.idx_t)(&centroids[0]),
		params.sp,
	); c != 0 {
		return nil, nil, getLastError()
	}

	return centroids, centroidDistances, nil
}

func (idx *faissIndex) ObtainKCentroidCardinalitiesFromIVFIndex(limit int, descending bool) (
	[]uint64, [][]float32, error) {
	if limit <= 0 {
		return nil, nil, nil
	}

	nlist := int(C.faiss_IndexIVF_nlist(idx.idx))
	if nlist == 0 {
		return nil, nil, nil
	}

	centroidCardinalities := make([]C.size_t, nlist)

	// Allocate a flat buffer for all centroids, then slice it per centroid
	d := idx.D()
	flatCentroids := make([]float32, nlist*d)

	// Call the C function to fill centroid vectors and cardinalities
	c := C.faiss_IndexIVF_get_centroids_and_cardinality(
		idx.idx,
		(*C.float)(&flatCentroids[0]),
		(*C.size_t)(&centroidCardinalities[0]),
		nil,
	)
	if c != 0 {
		return nil, nil, getLastError()
	}

	topIndices := getIndicesOfKCentroidCardinalities(
		centroidCardinalities,
		min(limit, nlist),
		descending)

	rvCardinalities := make([]uint64, len(topIndices))
	rvCentroids := make([][]float32, len(topIndices))

	for i, idx := range topIndices {
		rvCardinalities[i] = uint64(centroidCardinalities[idx])
		rvCentroids[i] = flatCentroids[idx*d : (idx+1)*d]
	}

	return rvCardinalities, rvCentroids, nil

}

func getIndicesOfKCentroidCardinalities(cardinalities []C.size_t, k int, descending bool) []int {
	n := len(cardinalities)
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	// Sort only the indices based on cardinality values
	sort.Slice(indices, func(i, j int) bool {
		if descending {
			return cardinalities[indices[i]] > cardinalities[indices[j]]
		}
		return cardinalities[indices[i]] < cardinalities[indices[j]]
	})
	if k >= n {
		return indices
	}

	return indices[:k]
}

func (idx *faissIndex) SearchClustersFromIVFIndex(eligibleCentroidIDs []int64, centroidDis []float32, centroidsToProbe int,
	x []float32, k int64, include Selector, params json.RawMessage) ([]float32, []int64, error) {
	// Applicable only to IVF indexes
	ivfPtr := C.faiss_IndexIVF_cast(idx.cPtr())
	if ivfPtr == nil {
		return nil, nil, fmt.Errorf("index is not of ivf type")
	}
	// If no include selector is provided, we have no results to return.
	// return an error indicating that the SearchClustersFromIVFIndex requires a valid selector.
	if include == nil {
		return nil, nil, fmt.Errorf("SearchClustersFromIVFIndex requires a valid include selector")
	}
	// create a temporary search params object to set nprobe, this will override
	// the nprobe and the nlist set at index time, this will allow the search to
	// probe only the clusters specified in eligibleCentroidIDs
	tempParams := &defaultSearchParamsIVF{
		// Nlist is set to the number of eligible centroids, which will override
		// the nlist set at index time.
		Nlist: len(eligibleCentroidIDs),
		// Have to override nprobe so that more clusters will be searched for this
		// query, if required.
		Nprobe: centroidsToProbe,
	}
	searchParams, err := NewSearchParams(idx, params, include, tempParams)
	if err != nil {
		return nil, nil, err
	}
	defer searchParams.Delete()

	n := len(x) / idx.D()

	distances := make([]float32, int64(n)*k)
	labels := make([]int64, int64(n)*k)
	// Adjust the slices to match the effective nprobe set in searchParams, as the input
	// parameters may have different nprobe value, which will be a hard override, over the
	// centroidsToProbe value passed to this function.
	// If the effective nprobe is greater than the length of eligibleCentroidIDs,
	// we limit it to the length of eligibleCentroidIDs.
	effectiveNprobe := min(getNProbeFromSearchParams(searchParams), int32(len(eligibleCentroidIDs)))
	eligibleCentroidIDs = eligibleCentroidIDs[:effectiveNprobe]
	centroidDis = centroidDis[:effectiveNprobe]

	if c := C.faiss_IndexIVF_search_preassigned_with_params(
		ivfPtr,
		(C.idx_t)(n),
		(*C.float)(&x[0]),
		(C.idx_t)(k),
		(*C.idx_t)(&eligibleCentroidIDs[0]),
		(*C.float)(&centroidDis[0]),
		(*C.float)(&distances[0]),
		(*C.idx_t)(&labels[0]),
		(C.int)(0),
		searchParams.sp,
	); c != 0 {
		return nil, nil, getLastError()
	}

	return distances, labels, nil
}

func (idx *faissIndex) AddWithIDs(x []float32, xids []int64) error {
	n := len(x) / idx.D()
	if c := C.faiss_Index_add_with_ids(
		idx.idx,
		C.idx_t(n),
		(*C.float)(&x[0]),
		(*C.idx_t)(&xids[0]),
	); c != 0 {
		return getLastError()
	}
	return nil
}

func (idx *faissIndex) Search(x []float32, k int64) (
	distances []float32, labels []int64, err error,
) {
	n := len(x) / idx.D()
	distances = make([]float32, int64(n)*k)
	labels = make([]int64, int64(n)*k)
	if c := C.faiss_Index_search(
		idx.idx,
		C.idx_t(n),
		(*C.float)(&x[0]),
		C.idx_t(k),
		(*C.float)(&distances[0]),
		(*C.idx_t)(&labels[0]),
	); c != 0 {
		err = getLastError()
	}

	return
}

// SearchWithoutIDs performs a search excluding the IDs specified in the exclude selector.
func (idx *faissIndex) SearchWithoutIDs(x []float32, k int64, exclude Selector, params json.RawMessage) (
	distances []float32, labels []int64, err error,
) {
	// If no exclude selector and no additional parameters are provided,
	// perform a standard search.
	if params == nil && exclude == nil {
		return idx.Search(x, k)
	}
	// Create search parameters with the exclude selector.
	searchParams, err := NewSearchParams(idx, params, exclude, nil)
	if err != nil {
		return nil, nil, err
	}
	// cleanup the searchParams after use
	defer searchParams.Delete()
	// Perform the search with the specified parameters.
	distances, labels, err = idx.searchWithParams(x, k, searchParams.sp)
	return
}

// SearchWithIDs performs a search including only the IDs specified in the include selector.
func (idx *faissIndex) SearchWithIDs(x []float32, k int64, include Selector, params json.RawMessage) (
	distances []float32, labels []int64, err error,
) {
	// If no include selector is provided, we have no results to return.
	// return an error indicating that the SearchWithIDs requires a valid selector.
	if include == nil {
		return nil, nil, fmt.Errorf("SearchWithIDs requires a valid include selector")
	}
	// Create search parameters with the include selector.
	searchParams, err := NewSearchParams(idx, params, include, nil)
	if err != nil {
		return nil, nil, err
	}
	// cleanup the searchParams after use
	defer searchParams.Delete()
	// Perform the search with the specified parameters.
	distances, labels, err = idx.searchWithParams(x, k, searchParams.sp)
	return
}

func (idx *faissIndex) Reconstruct(key int64) (recons []float32, err error) {
	rv := make([]float32, idx.D())
	if c := C.faiss_Index_reconstruct(
		idx.idx,
		C.idx_t(key),
		(*C.float)(&rv[0]),
	); c != 0 {
		err = getLastError()
	}

	return rv, err
}

func (idx *faissIndex) ReconstructBatch(keys []int64, recons []float32) ([]float32, error) {
	var err error
	n := int64(len(keys))
	if c := C.faiss_Index_reconstruct_batch(
		idx.idx,
		C.idx_t(n),
		(*C.idx_t)(&keys[0]),
		(*C.float)(&recons[0]),
	); c != 0 {
		err = getLastError()
	}

	return recons, err
}

func (i *IndexImpl) MergeFrom(other Index, add_id int64) error {
	if impl, ok := other.(*IndexImpl); ok {
		return i.Index.MergeFrom(impl.Index, add_id)
	}
	return fmt.Errorf("merge not support")
}

func (idx *faissIndex) MergeFrom(other Index, add_id int64) (err error) {
	otherIdx, ok := other.(*faissIndex)
	if !ok {
		return fmt.Errorf("merge api not supported")
	}

	if c := C.faiss_Index_merge_from(
		idx.idx,
		otherIdx.idx,
		(C.idx_t)(add_id),
	); c != 0 {
		err = getLastError()
	}

	return err
}

func (idx *faissIndex) RangeSearch(x []float32, radius float32) (
	*RangeSearchResult, error,
) {
	n := len(x) / idx.D()
	var rsr *C.FaissRangeSearchResult
	if c := C.faiss_RangeSearchResult_new(&rsr, C.idx_t(n)); c != 0 {
		return nil, getLastError()
	}
	if c := C.faiss_Index_range_search(
		idx.idx,
		C.idx_t(n),
		(*C.float)(&x[0]),
		C.float(radius),
		rsr,
	); c != 0 {
		return nil, getLastError()
	}
	return &RangeSearchResult{rsr}, nil
}

func (idx *faissIndex) DistCompute(queryData []float32, ids []int64) ([]float32, error) {
	distances := make([]float32, len(ids))
	if c := C.faiss_Index_dist_compute(idx.idx, (*C.float)(&queryData[0]),
		(*C.idx_t)(&ids[0]), (C.size_t)(len(ids)), (*C.float)(&distances[0])); c != 0 {
		return nil, getLastError()
	}

	return distances, nil
}

func (idx *faissIndex) Reset() error {
	if c := C.faiss_Index_reset(idx.idx); c != 0 {
		return getLastError()
	}
	return nil
}

func (idx *faissIndex) RemoveIDs(sel *IDSelector) (int, error) {
	var nRemoved C.size_t
	if c := C.faiss_Index_remove_ids(idx.idx, sel.sel, &nRemoved); c != 0 {
		return 0, getLastError()
	}
	return int(nRemoved), nil
}

func (idx *faissIndex) Close() {
	C.faiss_Index_free(idx.idx)
}

func (idx *faissIndex) searchWithParams(x []float32, k int64, searchParams *C.FaissSearchParameters) (
	distances []float32, labels []int64, err error,
) {
	n := len(x) / idx.D()
	distances = make([]float32, int64(n)*k)
	labels = make([]int64, int64(n)*k)

	if c := C.faiss_Index_search_with_params(
		idx.idx,
		C.idx_t(n),
		(*C.float)(&x[0]),
		C.idx_t(k),
		searchParams,
		(*C.float)(&distances[0]),
		(*C.idx_t)(&labels[0]),
	); c != 0 {
		err = getLastError()
	}

	return
}

// -----------------------------------------------------------------------------

// RangeSearchResult is the result of a range search.
type RangeSearchResult struct {
	rsr *C.FaissRangeSearchResult
}

// Nq returns the number of queries.
func (r *RangeSearchResult) Nq() int {
	return int(C.faiss_RangeSearchResult_nq(r.rsr))
}

// Lims returns a slice containing start and end indices for queries in the
// distances and labels slices returned by Labels.
func (r *RangeSearchResult) Lims() []int {
	var lims *C.size_t
	C.faiss_RangeSearchResult_lims(r.rsr, &lims)
	length := r.Nq() + 1
	return (*[1 << 30]int)(unsafe.Pointer(lims))[:length:length]
}

// Labels returns the unsorted IDs and respective distances for each query.
// The result for query i is labels[lims[i]:lims[i+1]].
func (r *RangeSearchResult) Labels() (labels []int64, distances []float32) {
	lims := r.Lims()
	length := lims[len(lims)-1]
	var clabels *C.idx_t
	var cdist *C.float
	C.faiss_RangeSearchResult_labels(r.rsr, &clabels, &cdist)
	labels = (*[1 << 30]int64)(unsafe.Pointer(clabels))[:length:length]
	distances = (*[1 << 30]float32)(unsafe.Pointer(cdist))[:length:length]
	return
}

// Delete frees the memory associated with r.
func (r *RangeSearchResult) Delete() {
	C.faiss_RangeSearchResult_free(r.rsr)
}

// IndexImpl is an abstract structure for an index.
type IndexImpl struct {
	Index
}

// IndexFactory builds a composite index.
// description is a comma-separated list of components.
func IndexFactory(d int, description string, metric int) (*IndexImpl, error) {
	cdesc := C.CString(description)
	defer C.free(unsafe.Pointer(cdesc))
	var idx faissIndex
	c := C.faiss_index_factory(&idx.idx, C.int(d), cdesc, C.FaissMetricType(metric))
	if c != 0 {
		return nil, getLastError()
	}
	return &IndexImpl{&idx}, nil
}

func SetOMPThreads(n uint) {
	C.faiss_set_omp_threads(C.uint(n))
}
