// SPDX-License-Identifier: Unlicense OR MIT

package headless

import (
	"errors"
	"unsafe"

	"gioui.org/gpu"
	_ "gioui.org/internal/cocoainit"
)

/*
#cgo CFLAGS: -Werror -Wno-deprecated-declarations -fmodules -fobjc-arc -x objective-c
#cgo LDFLAGS: -framework CoreGraphics

@import Metal;

static CFTypeRef createDevice(void) {
	@autoreleasepool {
		id dev = MTLCreateSystemDefaultDevice();
		return CFBridgingRetain(dev);
	}
}

static CFTypeRef newCommandQueue(CFTypeRef devRef) {
	@autoreleasepool {
		id<MTLDevice> dev = (__bridge id<MTLDevice>)devRef;
		return CFBridgingRetain([dev newCommandQueue]);
	}
}
*/
import "C"

type mtlContext struct {
	dev   C.CFTypeRef
	queue C.CFTypeRef
}

func newContext() (context, error) {
	dev := C.createDevice()
	if dev == 0 {
		return nil, errors.New("headless: failed to create Metal device")
	}
	queue := C.newCommandQueue(dev)
	if queue == 0 {
		C.CFRelease(dev)
		return nil, errors.New("headless: failed to create MTLQueue")
	}
	return &mtlContext{dev: dev, queue: queue}, nil
}

func (c *mtlContext) API() gpu.API {
	return gpu.Metal{
		Device:      unsafe.Pointer(c.dev),
		Queue:       unsafe.Pointer(c.queue),
		PixelFormat: int(C.MTLPixelFormatRGBA8Unorm_sRGB),
	}
}

func (c *mtlContext) MakeCurrent() error {
	return nil
}

func (c *mtlContext) ReleaseCurrent() {}

func (d *mtlContext) Release() {
	C.CFRelease(d.dev)
	C.CFRelease(d.queue)
	*d = mtlContext{}
}
