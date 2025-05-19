// main.go
package main

import (
    "log"
    "runtime"

    "github.com/veandco/go-sdl2/sdl"
    "github.com/go-gl/gl/v3.2-core/gl"
    imgui "github.com/inkyblackness/imgui-go/v4"
    "github.com/AEROSIQ/AEROSIQ-CLIENT/backends/platforms"
    "github.com/AEROSIQ/AEROSIQ-CLIENT/backends/renderers/opengl2"
)

func init() {
    // ImGui and SDL must run on the main OS thread
    runtime.LockOSThread()
}

func main() {
    // 1) Initialize SDL2 (all subsystems)
    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        log.Fatalf("sdl.Init failed: %v", err)
    }
    defer sdl.Quit()

    // 2) Request an OpenGL 2.1 context
    sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 2)
    sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 1)
    sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)

    const width, height = 800, 600

    // 3) Create the SDL window with OpenGL support
    window, err := sdl.CreateWindow(
        "Dear ImGui – SDL2 + Go",
        sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
        width, height,
        sdl.WINDOW_OPENGL|sdl.WINDOW_RESIZABLE,
    )
    if err != nil {
        log.Fatalf("CreateWindow failed: %v", err)
    }
    defer window.Destroy()

    // 4) Create the OpenGL context
    glCtx, err := window.GLCreateContext()
    if err != nil {
        log.Fatalf("GLCreateContext failed: %v", err)
    }
    defer sdl.GLDeleteContext(glCtx)

    // 5) Initialize OpenGL bindings
    if err := gl.Init(); err != nil {
        log.Fatalf("gl.Init failed: %v", err)
    }

    // 6) Set up ImGui context
    ctx := imgui.CreateContext(nil)
    defer ctx.Destroy()
    io := imgui.CurrentIO()
    io.SetDisplaySize(imgui.Vec2{X: float32(width), Y: float32(height)})

    // 7) Initialize your vendored SDL2 platform & OpenGL2 renderer
    platform := platforms.New(window)
    defer platform.Dispose()
    renderer := opengl2.New()
    defer renderer.Dispose()

    // 8) Main loop
    running := true
    for running {
        // a) Poll SDL events
        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch e := event.(type) {
            case *sdl.QuitEvent:
                running = false
            case *sdl.WindowEvent:
                if e.Event == sdl.WINDOWEVENT_CLOSE {
                    running = false
                }
            }
            platform.ProcessEvent(event)
        }

        // b) Start new ImGui frame
        imgui.NewFrame()
        platform.NewFrame()
        renderer.NewFrame()

        // c) Build your UI
        imgui.Begin("Hello, SDL2 + ImGui")
        imgui.Text("This is Dear ImGui running with SDL2 in Go.")
        imgui.End()

        // d) Render
        imgui.Render()
        renderer.Render(imgui.RenderedDrawData())

        // e) Swap buffers
        sdl.GL_SwapWindow(window)
    }
}
