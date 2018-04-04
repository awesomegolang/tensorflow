// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package tensorflow

import (
	"C"
	"fmt"
	"unsafe"
)

// #include <string.h>
// #include "tensorflow.h"
import "C"

var _ Tensor = &Float32Tensor{}

type Float32Tensor struct {
	dims []int
	data []float32
}

func NewFloat32Tensor(dims []int) *Float32Tensor {
	size := 1
	for _, dim := range dims {
		size *= dim
	}

	return &Float32Tensor{
		dims: dims,
		data: make([]float32, size),
	}
}

func (t *Float32Tensor) index(idx []int) (int, int) {
	if len(idx) >= len(t.dims) {
		panic(fmt.Sprintf("Trying to address using %d dimensions, only %d permitted",
			len(idx), len(t.dims)-1))
	}

	// Special case: the full array
	if len(idx) == 0 {
		dimSize := 1
		for _, dim := range t.dims {
			dimSize *= dim
		}

		return 0, dimSize
	}

	startIdx := 0
	var dimSize int
	for i, idx := range idx {
		dimSize = 1
		for _, dim := range t.dims[i+1:] {
			dimSize *= dim
		}

		startIdx += idx * dimSize
	}

	return startIdx, dimSize
}

func (t *Float32Tensor) Get(idx []int) []float32 {
	startIdx, dimSize := t.index(idx)

	return t.data[startIdx : startIdx+dimSize]
}

func (t *Float32Tensor) Assign(idx []int, data []float32) {
	startIdx, _ := t.index(idx)

	if len(idx) >= len(t.dims) {
		panic(fmt.Sprintf("Trying to address using %d dimensions, only %d permitted",
			len(idx), len(t.dims)-1))
	}

	// TODO: bounds checks?

	copy(t.data[startIdx:], data)
}

func (t *Float32Tensor) Fill(idx []int, v float32) {
	startIdx, dimSize := t.index(idx)

	for idx := startIdx; idx < startIdx+dimSize; idx++ {
		t.data[idx] = v
	}
}

func (t *Float32Tensor) ToNative() *NativeTensor {
	// TF_NewTensor copies dims, does not take ownership.
	llDims := make([]C.longlong, len(t.dims))
	for idx, val := range t.dims {
		llDims[idx] = C.longlong(val)
	}

	dataLen := C.size_t(len(t.data)) * C.size_t(unsafe.Sizeof(t.data[0]))

	// Allocate new memory, rather than using the Go slice backing array,
	// since we cannot fullfil the alignment preferences.
	cTensor := C.TF_AllocateTensor(C.TF_FLOAT, (*C.int64_t)(unsafe.Pointer(&llDims[0])),
		C.int(len(llDims)), dataLen)
	cData := C.TF_TensorData(cTensor)

	C.memcpy(cData, unsafe.Pointer(&t.data[0]), dataLen)

	return &NativeTensor{
		inner: cTensor,
	}
}

func adoptfloat32Tensor(ct *C.TF_Tensor) *Float32Tensor {
	dims := C.TF_NumDims(ct)
	shape := make([]int, dims)
	size := uint(1)
	for i := C.int(0); i < dims; i++ {
		shape[i] = int(C.TF_Dim(ct, i))
		size *= uint(shape[i])
	}

	bs := C.TF_TensorByteSize(ct)
	var valForSize float32
	if uint(bs)/uint(unsafe.Sizeof(valForSize)) != size {
		panic("Expected tensor size does not correspond to the actual tensor size")
	}

	data := make([]float32, size)
	C.memcpy(unsafe.Pointer(&data[0]), C.TF_TensorData(ct), bs)

	return &Float32Tensor{
		dims: shape,
		data: data,
	}
}
