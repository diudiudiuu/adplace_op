# Wails项目构建Makefile

.PHONY: run build-win build-mac fix-mac

# 开发运行
run:
	@echo "启动开发服务器..."
	wails dev

# 构建Windows版本
build-win:
	@echo "构建Windows版本..."
	wails build -platform windows/amd64 -o adsplat-windows.exe

# 构建macOS版本 (Apple Silicon M芯片)
build-mac:
	@echo "构建macOS版本 (Apple Silicon M芯片)..."
	wails build -platform darwin/arm64 -o adsplat-macos-arm.app
	@echo "移除隔离属性..."
	xattr -rd com.apple.quarantine build/bin/adsplat-macos-arm.app 2>/dev/null || true

# 修复macOS应用权限问题
fix-mac:
	@echo "修复macOS应用权限..."
	xattr -rd com.apple.quarantine build/bin/adsplat-macos-arm.app 2>/dev/null || true
	@echo "权限修复完成，现在可以运行应用了"