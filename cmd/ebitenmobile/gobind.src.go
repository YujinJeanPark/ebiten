// Code generated by file2byteslice. DO NOT EDIT.

package main

var gobindsrc = []byte("// Copyright 2019 The Ebiten Authors\n//\n// Licensed under the Apache License, Version 2.0 (the \"License\");\n// you may not use this file except in compliance with the License.\n// You may obtain a copy of the License at\n//\n//     http://www.apache.org/licenses/LICENSE-2.0\n//\n// Unless required by applicable law or agreed to in writing, software\n// distributed under the License is distributed on an \"AS IS\" BASIS,\n// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n// See the License for the specific language governing permissions and\n// limitations under the License.\n\n// +build ebitenmobilegobind\n\n// gobind is a wrapper of the original gobind. This command adds extra files like a view controller.\npackage main\n\nimport (\n\t\"flag\"\n\t\"fmt\"\n\t\"io/ioutil\"\n\t\"log\"\n\t\"os\"\n\t\"os/exec\"\n\t\"path/filepath\"\n\t\"strings\"\n\n\t\"golang.org/x/tools/go/packages\"\n)\n\nvar (\n\tlang          = flag.String(\"lang\", \"\", \"\")\n\toutdir        = flag.String(\"outdir\", \"\", \"\")\n\tjavaPkg       = flag.String(\"javapkg\", \"\", \"\")\n\tprefix        = flag.String(\"prefix\", \"\", \"\")\n\tbootclasspath = flag.String(\"bootclasspath\", \"\", \"\")\n\tclasspath     = flag.String(\"classpath\", \"\", \"\")\n\ttags          = flag.String(\"tags\", \"\", \"\")\n)\n\nvar usage = `The Gobind tool generates Java language bindings for Go.\n\nFor usage details, see doc.go.`\n\nfunc main() {\n\tflag.Parse()\n\tif err := run(); err != nil {\n\t\tlog.Fatal(err)\n\t}\n}\n\nfunc invokeOriginalGobind(lang string) (pkgName string, err error) {\n\tcmd := exec.Command(\"gobind-original\", os.Args[1:]...)\n\tcmd.Stdout = os.Stdout\n\tcmd.Stderr = os.Stderr\n\tif err := cmd.Run(); err != nil {\n\t\treturn \"\", err\n\t}\n\n\tcfgtags := strings.Join(strings.Split(*tags, \",\"), \" \")\n\tcfg := &packages.Config{}\n\tswitch lang {\n\tcase \"java\":\n\t\tcfg.Env = append(os.Environ(), \"GOOS=android\")\n\tcase \"objc\":\n\t\tcfg.Env = append(os.Environ(), \"GOOS=darwin\")\n\t\tif cfgtags != \"\" {\n\t\t\tcfgtags += \" \"\n\t\t}\n\t\tcfgtags += \"ios\"\n\t}\n\tcfg.BuildFlags = []string{\"-tags\", cfgtags}\n\tpkgs, err := packages.Load(cfg, flag.Args()[0])\n\tif err != nil {\n\t\treturn \"\", err\n\t}\n\treturn pkgs[0].Name, nil\n}\n\nfunc forceGL() bool {\n\tfor _, tag := range strings.Split(*tags, \",\") {\n\t\tif tag == \"ebitengl\" {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false\n}\n\nfunc run() error {\n\twriteFile := func(filename string, content string) error {\n\t\tif err := ioutil.WriteFile(filepath.Join(*outdir, filename), []byte(content), 0644); err != nil {\n\t\t\treturn err\n\t\t}\n\t\treturn nil\n\t}\n\n\t// Add additional files.\n\tlangs := strings.Split(*lang, \",\")\n\tfor _, lang := range langs {\n\t\tpkgName, err := invokeOriginalGobind(lang)\n\t\tif err != nil {\n\t\t\treturn err\n\t\t}\n\t\tprefixLower := *prefix + pkgName\n\t\tprefixUpper := strings.Title(*prefix) + strings.Title(pkgName)\n\t\treplacePrefixes := func(content string) string {\n\t\t\tcontent = strings.ReplaceAll(content, \"{{.PrefixUpper}}\", prefixUpper)\n\t\t\tcontent = strings.ReplaceAll(content, \"{{.PrefixLower}}\", prefixLower)\n\t\t\tcontent = strings.ReplaceAll(content, \"{{.JavaPkg}}\", *javaPkg)\n\n\t\t\tf := \"0\"\n\t\t\tif forceGL() {\n\t\t\t\tf = \"1\"\n\t\t\t}\n\t\t\tcontent = strings.ReplaceAll(content, \"{{.ForceGL}}\", f)\n\t\t\treturn content\n\t\t}\n\n\t\tswitch lang {\n\t\tcase \"objc\":\n\t\t\t// iOS\n\t\t\tif err := writeFile(filepath.Join(\"src\", \"gobind\", prefixLower+\"ebitenviewcontroller_ios.m\"), replacePrefixes(objcM)); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\tcase \"java\":\n\t\t\t// Android\n\t\t\tdir := filepath.Join(strings.Split(*javaPkg, \".\")...)\n\t\t\tdir = filepath.Join(dir, prefixLower)\n\t\t\tif err := writeFile(filepath.Join(\"java\", dir, \"EbitenView.java\"), replacePrefixes(viewJava)); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\t\tif err := writeFile(filepath.Join(\"java\", dir, \"EbitenSurfaceView.java\"), replacePrefixes(surfaceViewJava)); err != nil {\n\t\t\t\treturn err\n\t\t\t}\n\t\tcase \"go\":\n\t\t\t// Do nothing.\n\t\tdefault:\n\t\t\tpanic(fmt.Sprintf(\"unsupported language: %s\", lang))\n\t\t}\n\t}\n\n\treturn nil\n}\n\nconst objcM = `// Code generated by ebitenmobile. DO NOT EDIT.\n\n// +build ios\n\n#import <TargetConditionals.h>\n\n#if TARGET_IPHONE_SIMULATOR || {{.ForceGL}}\n#define EBITEN_METAL 0\n#else\n#define EBITEN_METAL 1\n#endif\n\n#import <stdint.h>\n#import <UIKit/UIKit.h>\n#import <GLKit/GLkit.h>\n\n#import \"Ebitenmobileview.objc.h\"\n\n@interface {{.PrefixUpper}}EbitenViewController : UIViewController\n@end\n\n@implementation {{.PrefixUpper}}EbitenViewController {\n  UIView*  metalView_;\n  GLKView* glkView_;\n  bool     started_;\n  bool     active_;\n  bool     error_;\n}\n\n- (UIView*)metalView {\n  if (!metalView_) {\n    metalView_ = [[UIView alloc] init];\n    metalView_.multipleTouchEnabled = YES;\n  }\n  return metalView_;\n}\n\n- (GLKView*)glkView {\n  if (!glkView_) {\n    glkView_ = [[GLKView alloc] init];\n    glkView_.multipleTouchEnabled = YES;\n  }\n  return glkView_;\n}\n\n- (void)viewDidLoad {\n  [super viewDidLoad];\n\n  if (!started_) {\n    @synchronized(self) {\n      active_ = true;\n    }\n    started_ = true;\n  }\n\n#if EBITEN_METAL\n  [self.view addSubview: self.metalView];\n  EbitenmobileviewSetUIView((uintptr_t)(self.metalView));\n#else\n  self.glkView.delegate = (id<GLKViewDelegate>)(self);\n  [self.view addSubview: self.glkView];\n\n  EAGLContext *context = [[EAGLContext alloc] initWithAPI:kEAGLRenderingAPIOpenGLES2];\n  [self glkView].context = context;\n\t\n  [EAGLContext setCurrentContext:context];\n#endif\n\n  CADisplayLink *displayLink = [CADisplayLink displayLinkWithTarget:self selector:@selector(drawFrame)];\n  [displayLink addToRunLoop:[NSRunLoop currentRunLoop] forMode:NSDefaultRunLoopMode];\n}\n\n- (void)viewWillLayoutSubviews {\n  CGRect viewRect = [[self view] frame];\n#if EBITEN_METAL\n  [[self metalView] setFrame:viewRect];\n#else\n  [[self glkView] setFrame:viewRect];\n#endif\n}\n\n- (void)viewDidLayoutSubviews {\n  [super viewDidLayoutSubviews];\n  CGRect viewRect = [[self view] frame];\n\n  EbitenmobileviewLayout(viewRect.size.width, viewRect.size.height);\n}\n\n- (void)didReceiveMemoryWarning {\n  [super didReceiveMemoryWarning];\n  // Dispose of any resources that can be recreated.\n  // TODO: Notify this to Go world?\n}\n\n- (void)drawFrame{\n  @synchronized(self) {\n    if (!active_) {\n      return;\n    }\n\n#if EBITEN_METAL\n    [self updateEbiten];\n#else\n    [[self glkView] setNeedsDisplay];\n#endif\n  }\n}\n\n- (void)glkView:(GLKView*)view drawInRect:(CGRect)rect {\n  @synchronized(self) {\n    [self updateEbiten];\n  }\n}\n\n- (void)updateEbiten {\n  if (error_) {\n    return;\n  }\n  NSError* err = nil;\n  EbitenmobileviewUpdate(&err);\n  if (err != nil) {\n    [self performSelectorOnMainThread:@selector(onErrorOnGameUpdate:)\n                           withObject:err\n                        waitUntilDone:NO];\n    error_ = true;\n  }\n}\n\n- (void)onErrorOnGameUpdate:(NSError*)err {\n  NSLog(@\"Error: %@\", err);\n}\n\n- (void)updateTouches:(NSSet*)touches {\n  for (UITouch* touch in touches) {\n#if EBITEN_METAL\n    if (touch.view != [self metalView]) {\n      continue;\n    }\n#else\n    if (touch.view != [self glkView]) {\n      continue;\n    }\n#endif\n    CGPoint location = [touch locationInView:touch.view];\n    EbitenmobileviewUpdateTouchesOnIOS(touch.phase, (uintptr_t)touch, location.x, location.y);\n  }\n}\n\n- (void)touchesBegan:(NSSet*)touches withEvent:(UIEvent*)event {\n  [self updateTouches:touches];\n}\n\n- (void)touchesMoved:(NSSet*)touches withEvent:(UIEvent*)event {\n  [self updateTouches:touches];\n}\n\n- (void)touchesEnded:(NSSet*)touches withEvent:(UIEvent*)event {\n  [self updateTouches:touches];\n}\n\n- (void)touchesCancelled:(NSSet*)touches withEvent:(UIEvent*)event {\n  [self updateTouches:touches];\n}\n\n- (void)suspendGame {\n  NSAssert(started_, @\"suspendGame must not be called before viewDidLoad is called\");\n\n  @synchronized(self) {\n    active_ = false;\n    EbitenmobileviewSuspend();\n  }\n}\n\n- (void)resumeGame {\n  NSAssert(started_, @\"resumeGame must not be called before viewDidLoad is called\");\n\n  @synchronized(self) {\n    active_ = true;\n    EbitenmobileviewResume();\n  }\n}\n\n@end\n`\n\nconst viewJava = `// Code generated by ebitenmobile. DO NOT EDIT.\n\npackage {{.JavaPkg}}.{{.PrefixLower}};\n\nimport android.content.Context;\nimport android.hardware.input.InputManager;\nimport android.os.Handler;\nimport android.os.Looper;\nimport android.util.AttributeSet;\nimport android.util.Log;\nimport android.view.KeyEvent;\nimport android.view.InputDevice;\nimport android.view.MotionEvent;\nimport android.view.ViewGroup;\n\nimport {{.JavaPkg}}.ebitenmobileview.Ebitenmobileview;\n\npublic class EbitenView extends ViewGroup implements InputManager.InputDeviceListener {\n    private double getDeviceScale() {\n        if (this.deviceScale == 0.0) {\n            this.deviceScale = getResources().getDisplayMetrics().density;\n        }\n        return this.deviceScale;\n    }\n\n    private double pxToDp(double x) {\n        return x / getDeviceScale();\n    }\n\n    private double deviceScale = 0.0;\n\n    public EbitenView(Context context) {\n        super(context);\n        initialize(context);\n    }\n\n    public EbitenView(Context context, AttributeSet attrs) {\n        super(context, attrs);\n        initialize(context);\n    }\n\n    private void initialize(Context context) {\n        this.ebitenSurfaceView = new EbitenSurfaceView(getContext());\n        LayoutParams params = new LayoutParams(LayoutParams.MATCH_PARENT, LayoutParams.MATCH_PARENT);\n        addView(this.ebitenSurfaceView, params);\n\n        this.inputManager = (InputManager)context.getSystemService(Context.INPUT_SERVICE);\n        this.inputManager.registerInputDeviceListener(this, null);\n        for (int id : this.inputManager.getInputDeviceIds()) {\n            this.onInputDeviceAdded(id);\n        }\n    }\n\n    @Override\n    protected void onLayout(boolean changed, int left, int top, int right, int bottom) {\n        this.ebitenSurfaceView.layout(0, 0, right - left, bottom - top);\n        double widthInDp = pxToDp(right - left);\n        double heightInDp = pxToDp(bottom - top);\n        Ebitenmobileview.layout(widthInDp, heightInDp);\n    }\n\n    @Override\n    public boolean onKeyDown(int keyCode, KeyEvent event) {\n        Ebitenmobileview.onKeyDownOnAndroid(keyCode, event.getUnicodeChar(), event.getSource(), event.getDeviceId());\n        return true;\n    }\n\n    @Override\n    public boolean onKeyUp(int keyCode, KeyEvent event) {\n        Ebitenmobileview.onKeyUpOnAndroid(keyCode, event.getSource(), event.getDeviceId());\n        return true;\n    }\n\n    @Override\n    public boolean onTouchEvent(MotionEvent e) {\n        for (int i = 0; i < e.getPointerCount(); i++) {\n            int id = e.getPointerId(i);\n            int x = (int)e.getX(i);\n            int y = (int)e.getY(i);\n            Ebitenmobileview.updateTouchesOnAndroid(e.getActionMasked(), id, (int)pxToDp(x), (int)pxToDp(y));\n        }\n        return true;\n    }\n\n    // The order must be the same as mobile/ebitenmobileview/input_android.go.\n    static int[] gamepadButtons = {\n        KeyEvent.KEYCODE_BUTTON_A,\n        KeyEvent.KEYCODE_BUTTON_B,\n        KeyEvent.KEYCODE_BUTTON_C,\n        KeyEvent.KEYCODE_BUTTON_X,\n        KeyEvent.KEYCODE_BUTTON_Y,\n        KeyEvent.KEYCODE_BUTTON_Z,\n        KeyEvent.KEYCODE_BUTTON_L1,\n        KeyEvent.KEYCODE_BUTTON_R1,\n        KeyEvent.KEYCODE_BUTTON_L2,\n        KeyEvent.KEYCODE_BUTTON_R2,\n        KeyEvent.KEYCODE_BUTTON_THUMBL,\n        KeyEvent.KEYCODE_BUTTON_THUMBR,\n        KeyEvent.KEYCODE_BUTTON_START,\n        KeyEvent.KEYCODE_BUTTON_SELECT,\n        KeyEvent.KEYCODE_BUTTON_MODE,\n        KeyEvent.KEYCODE_BUTTON_1,\n        KeyEvent.KEYCODE_BUTTON_2,\n        KeyEvent.KEYCODE_BUTTON_3,\n        KeyEvent.KEYCODE_BUTTON_4,\n        KeyEvent.KEYCODE_BUTTON_5,\n        KeyEvent.KEYCODE_BUTTON_6,\n        KeyEvent.KEYCODE_BUTTON_7,\n        KeyEvent.KEYCODE_BUTTON_8,\n        KeyEvent.KEYCODE_BUTTON_9,\n        KeyEvent.KEYCODE_BUTTON_10,\n        KeyEvent.KEYCODE_BUTTON_11,\n        KeyEvent.KEYCODE_BUTTON_12,\n        KeyEvent.KEYCODE_BUTTON_13,\n        KeyEvent.KEYCODE_BUTTON_14,\n        KeyEvent.KEYCODE_BUTTON_15,\n        KeyEvent.KEYCODE_BUTTON_16,\n    };\n\n    // The order must be the same as mobile/ebitenmobileview/input_android.go.\n    static int[] axes = {\n        MotionEvent.AXIS_X,\n        MotionEvent.AXIS_Y,\n        MotionEvent.AXIS_Z,\n        MotionEvent.AXIS_RX,\n        MotionEvent.AXIS_RY,\n        MotionEvent.AXIS_RZ,\n        MotionEvent.AXIS_HAT_X,\n        MotionEvent.AXIS_HAT_Y,\n        MotionEvent.AXIS_LTRIGGER,\n        MotionEvent.AXIS_RTRIGGER,\n        MotionEvent.AXIS_THROTTLE,\n        MotionEvent.AXIS_RUDDER,\n        MotionEvent.AXIS_WHEEL,\n        MotionEvent.AXIS_GAS,\n        MotionEvent.AXIS_BRAKE,\n        MotionEvent.AXIS_GENERIC_1,\n        MotionEvent.AXIS_GENERIC_2,\n        MotionEvent.AXIS_GENERIC_3,\n        MotionEvent.AXIS_GENERIC_4,\n        MotionEvent.AXIS_GENERIC_5,\n        MotionEvent.AXIS_GENERIC_6,\n        MotionEvent.AXIS_GENERIC_7,\n        MotionEvent.AXIS_GENERIC_8,\n        MotionEvent.AXIS_GENERIC_9,\n        MotionEvent.AXIS_GENERIC_10,\n        MotionEvent.AXIS_GENERIC_11,\n        MotionEvent.AXIS_GENERIC_12,\n        MotionEvent.AXIS_GENERIC_13,\n        MotionEvent.AXIS_GENERIC_14,\n        MotionEvent.AXIS_GENERIC_15,\n        MotionEvent.AXIS_GENERIC_16,\n    };\n\n    @Override\n    public boolean onGenericMotionEvent(MotionEvent event) {\n        if ((event.getSource() & InputDevice.SOURCE_JOYSTICK) != InputDevice.SOURCE_JOYSTICK) {\n            return super.onGenericMotionEvent(event);\n        }\n        if (event.getAction() != MotionEvent.ACTION_MOVE) {\n            return super.onGenericMotionEvent(event);\n        }\n        InputDevice inputDevice = this.inputManager.getInputDevice(event.getDeviceId());\n        for (int axis : axes) {\n            InputDevice.MotionRange motionRange = inputDevice.getMotionRange(axis, event.getSource());\n            float value = 0.0f;\n            if (motionRange != null) {\n                value = event.getAxisValue(axis);\n                if (Math.abs(value) <= motionRange.getFlat()) {\n                    value = 0.0f;\n                }\n            }\n            Ebitenmobileview.onGamepadAxesChanged(event.getDeviceId(), axis, value);\n        }\n        return true;\n    }\n\n    @Override\n    public void onInputDeviceAdded(int deviceId) {\n        InputDevice inputDevice = this.inputManager.getInputDevice(deviceId);\n        int sources = inputDevice.getSources();\n        if ((sources & InputDevice.SOURCE_GAMEPAD) != InputDevice.SOURCE_GAMEPAD &&\n            (sources & InputDevice.SOURCE_JOYSTICK) != InputDevice.SOURCE_JOYSTICK) {\n            return;\n        }\n\n        boolean[] keyExistences = inputDevice.hasKeys(gamepadButtons);\n        int buttonNum = gamepadButtons.length - 1;\n        for (int i = gamepadButtons.length - 1; i >= 0; i--) {\n            if (keyExistences[i]) {\n                break;\n            }\n            buttonNum--;\n        }\n\n        int axisNum = axes.length - 1;\n        for (int i = axes.length - 1; i >= 0; i--) {\n            if (inputDevice.getMotionRange(axes[i], InputDevice.SOURCE_JOYSTICK) != null) {\n                break;\n            }\n            axisNum--;\n        }\n\n        String descriptor = inputDevice.getDescriptor();\n        int vendorId = inputDevice.getVendorId();\n        int productId = inputDevice.getProductId();\n\n        // These values are required to calculate SDL's GUID.\n        int buttonMask = getButtonMask(inputDevice);\n        int axisMask = getAxisMask(inputDevice);\n\n        Ebitenmobileview.onGamepadAdded(deviceId, inputDevice.getName(), buttonNum, axisNum, descriptor, vendorId, productId, buttonMask, axisMask);\n    }\n\n    // The implementation is copied from SDL:\n    // https://hg.libsdl.org/SDL/file/bc90ce38f1e2/android-project/app/src/main/java/org/libsdl/app/SDLControllerManager.java#l308\n    private int getButtonMask(InputDevice joystickDevice) {\n        int button_mask = 0;\n        int[] keys = new int[] {\n            KeyEvent.KEYCODE_BUTTON_A,\n            KeyEvent.KEYCODE_BUTTON_B,\n            KeyEvent.KEYCODE_BUTTON_X,\n            KeyEvent.KEYCODE_BUTTON_Y,\n            KeyEvent.KEYCODE_BACK,\n            KeyEvent.KEYCODE_BUTTON_MODE,\n            KeyEvent.KEYCODE_BUTTON_START,\n            KeyEvent.KEYCODE_BUTTON_THUMBL,\n            KeyEvent.KEYCODE_BUTTON_THUMBR,\n            KeyEvent.KEYCODE_BUTTON_L1,\n            KeyEvent.KEYCODE_BUTTON_R1,\n            KeyEvent.KEYCODE_DPAD_UP,\n            KeyEvent.KEYCODE_DPAD_DOWN,\n            KeyEvent.KEYCODE_DPAD_LEFT,\n            KeyEvent.KEYCODE_DPAD_RIGHT,\n            KeyEvent.KEYCODE_BUTTON_SELECT,\n            KeyEvent.KEYCODE_DPAD_CENTER,\n\n            // These don't map into any SDL controller buttons directly\n            KeyEvent.KEYCODE_BUTTON_L2,\n            KeyEvent.KEYCODE_BUTTON_R2,\n            KeyEvent.KEYCODE_BUTTON_C,\n            KeyEvent.KEYCODE_BUTTON_Z,\n            KeyEvent.KEYCODE_BUTTON_1,\n            KeyEvent.KEYCODE_BUTTON_2,\n            KeyEvent.KEYCODE_BUTTON_3,\n            KeyEvent.KEYCODE_BUTTON_4,\n            KeyEvent.KEYCODE_BUTTON_5,\n            KeyEvent.KEYCODE_BUTTON_6,\n            KeyEvent.KEYCODE_BUTTON_7,\n            KeyEvent.KEYCODE_BUTTON_8,\n            KeyEvent.KEYCODE_BUTTON_9,\n            KeyEvent.KEYCODE_BUTTON_10,\n            KeyEvent.KEYCODE_BUTTON_11,\n            KeyEvent.KEYCODE_BUTTON_12,\n            KeyEvent.KEYCODE_BUTTON_13,\n            KeyEvent.KEYCODE_BUTTON_14,\n            KeyEvent.KEYCODE_BUTTON_15,\n            KeyEvent.KEYCODE_BUTTON_16,\n        };\n        int[] masks = new int[] {\n            (1 << 0),   // A -> A\n            (1 << 1),   // B -> B\n            (1 << 2),   // X -> X\n            (1 << 3),   // Y -> Y\n            (1 << 4),   // BACK -> BACK\n            (1 << 5),   // MODE -> GUIDE\n            (1 << 6),   // START -> START\n            (1 << 7),   // THUMBL -> LEFTSTICK\n            (1 << 8),   // THUMBR -> RIGHTSTICK\n            (1 << 9),   // L1 -> LEFTSHOULDER\n            (1 << 10),  // R1 -> RIGHTSHOULDER\n            (1 << 11),  // DPAD_UP -> DPAD_UP\n            (1 << 12),  // DPAD_DOWN -> DPAD_DOWN\n            (1 << 13),  // DPAD_LEFT -> DPAD_LEFT\n            (1 << 14),  // DPAD_RIGHT -> DPAD_RIGHT\n            (1 << 4),   // SELECT -> BACK\n            (1 << 0),   // DPAD_CENTER -> A\n            (1 << 15),  // L2 -> ??\n            (1 << 16),  // R2 -> ??\n            (1 << 17),  // C -> ??\n            (1 << 18),  // Z -> ??\n            (1 << 20),  // 1 -> ??\n            (1 << 21),  // 2 -> ??\n            (1 << 22),  // 3 -> ??\n            (1 << 23),  // 4 -> ??\n            (1 << 24),  // 5 -> ??\n            (1 << 25),  // 6 -> ??\n            (1 << 26),  // 7 -> ??\n            (1 << 27),  // 8 -> ??\n            (1 << 28),  // 9 -> ??\n            (1 << 29),  // 10 -> ??\n            (1 << 30),  // 11 -> ??\n            (1 << 31),  // 12 -> ??\n            // We're out of room...\n            0xFFFFFFFF,  // 13 -> ??\n            0xFFFFFFFF,  // 14 -> ??\n            0xFFFFFFFF,  // 15 -> ??\n            0xFFFFFFFF,  // 16 -> ??\n        };\n        boolean[] has_keys = joystickDevice.hasKeys(keys);\n        for (int i = 0; i < keys.length; ++i) {\n            if (has_keys[i]) {\n                button_mask |= masks[i];\n            }\n        }\n        return button_mask;\n    }\n\n    private int getAxisMask(InputDevice joystickDevice) {\n        final int SDL_CONTROLLER_AXIS_LEFTX = 0;\n        final int SDL_CONTROLLER_AXIS_LEFTY = 1;\n        final int SDL_CONTROLLER_AXIS_RIGHTX = 2;\n        final int SDL_CONTROLLER_AXIS_RIGHTY = 3;\n        final int SDL_CONTROLLER_AXIS_TRIGGERLEFT = 4;\n        final int SDL_CONTROLLER_AXIS_TRIGGERRIGHT = 5;\n\n        int naxes = 0;\n        for (InputDevice.MotionRange range : joystickDevice.getMotionRanges()) {\n            if ((range.getSource() & InputDevice.SOURCE_CLASS_JOYSTICK) != 0) {\n                if (range.getAxis() != MotionEvent.AXIS_HAT_X && range.getAxis() != MotionEvent.AXIS_HAT_Y) {\n                    naxes++;\n                }\n            }\n        }\n        // The variable is_accelerometer seems always false, then skip the checking:\n        // https://hg.libsdl.org/SDL/file/bc90ce38f1e2/android-project/app/src/main/java/org/libsdl/app/SDLControllerManager.java#l207\n        int axisMask = 0;\n        if (naxes >= 2) {\n            axisMask |= ((1 << SDL_CONTROLLER_AXIS_LEFTX) | (1 << SDL_CONTROLLER_AXIS_LEFTY));\n        }\n        if (naxes >= 4) {\n            axisMask |= ((1 << SDL_CONTROLLER_AXIS_RIGHTX) | (1 << SDL_CONTROLLER_AXIS_RIGHTY));\n        }\n        if (naxes >= 6) {\n            axisMask |= ((1 << SDL_CONTROLLER_AXIS_TRIGGERLEFT) | (1 << SDL_CONTROLLER_AXIS_TRIGGERRIGHT));\n        }\n        return axisMask;\n    }\n\n    @Override\n    public void onInputDeviceChanged(int deviceId) {\n        // Do nothing.\n    }\n\n    @Override\n    public void onInputDeviceRemoved(int deviceId) {\n        // Do not call inputManager.getInputDevice(), which returns null (#1185).\n        Ebitenmobileview.onInputDeviceRemoved(deviceId);\n    }\n\n    // suspendGame suspends the game.\n    // It is recommended to call this when the application is being suspended e.g.,\n    // Activity's onPause is called.\n    public void suspendGame() {\n        this.inputManager.unregisterInputDeviceListener(this);\n        this.ebitenSurfaceView.onPause();\n        Ebitenmobileview.suspend();\n    }\n\n    // resumeGame resumes the game.\n    // It is recommended to call this when the application is being resumed e.g.,\n    // Activity's onResume is called.\n    public void resumeGame() {\n        this.inputManager.registerInputDeviceListener(this, null);\n        this.ebitenSurfaceView.onResume();\n        Ebitenmobileview.resume();\n    }\n\n    // onErrorOnGameUpdate is called on the main thread when an error happens when updating a game.\n    // You can define your own error handler, e.g., using Crashlytics, by overwriting this method.\n    protected void onErrorOnGameUpdate(Exception e) {\n        Log.e(\"Go\", e.toString());\n    }\n\n    private EbitenSurfaceView ebitenSurfaceView;\n    private InputManager inputManager;\n}\n`\n\nconst surfaceViewJava = `// Code generated by ebitenmobile. DO NOT EDIT.\n\npackage {{.JavaPkg}}.{{.PrefixLower}};\n\nimport android.content.Context;\nimport android.opengl.GLSurfaceView;\nimport android.os.Handler;\nimport android.os.Looper;\nimport android.util.AttributeSet;\nimport android.util.Log;\n\nimport javax.microedition.khronos.egl.EGLConfig;\nimport javax.microedition.khronos.opengles.GL10;\n\nimport {{.JavaPkg}}.ebitenmobileview.Ebitenmobileview;\nimport {{.JavaPkg}}.{{.PrefixLower}}.EbitenView;\n\nclass EbitenSurfaceView extends GLSurfaceView {\n\n    private class EbitenRenderer implements GLSurfaceView.Renderer {\n\n        private boolean errored_ = false;\n\n        @Override\n        public void onDrawFrame(GL10 gl) {\n            if (errored_) {\n                return;\n            }\n            try {\n                Ebitenmobileview.update();\n            } catch (final Exception e) {\n                new Handler(Looper.getMainLooper()).post(new Runnable() {\n                    @Override\n                    public void run() {\n                        onErrorOnGameUpdate(e);\n                    }\n                });\n                errored_ = true;\n            }\n        }\n\n        @Override\n        public void onSurfaceCreated(GL10 gl, EGLConfig config) {\n            Ebitenmobileview.onContextLost();\n        }\n\n        @Override\n        public void onSurfaceChanged(GL10 gl, int width, int height) {\n        }\n    }\n\n    public EbitenSurfaceView(Context context) {\n        super(context);\n        initialize();\n    }\n\n    public EbitenSurfaceView(Context context, AttributeSet attrs) {\n        super(context, attrs);\n        initialize();\n    }\n\n    private void initialize() {\n        setEGLContextClientVersion(2);\n        setEGLConfigChooser(8, 8, 8, 8, 0, 0);\n        setRenderer(new EbitenRenderer());\n    }\n\n    private void onErrorOnGameUpdate(Exception e) {\n        ((EbitenView)getParent()).onErrorOnGameUpdate(e);\n    }\n}\n`\n")
