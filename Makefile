ifdef OS
	QMAKE_RELEASE = release
	QTBINDING_DLL = standalone/qt/build/release/qtbinding.dll
	EXENAME = rxgui.exe
	DLLNAME = rxgui.dll
else
	QMAKE_RELEASE = .
	QTBINDING_DLL = standalone/qt/build/libqtbinding*
	EXENAME = rxgui
	DLLNAME = rxgui.so
endif

ifeq ($(QTBINDING_ASAN),1)
	ASAN_ENABLED = 1
	ASAN_FLAG = -asan
else
	ASAN_ENABLED = 0
	ASAN_FLAG =
endif

.PHONY: check qt qt-clean crash-report naive-debugger tools deps interpreter release
default: interpreter

check:
	@echo -e '\033[1mChecking for Qt...\033[0m'
	qmake -v
	@echo -e '\033[1mChecking for Go...\033[0m'
	go version

qt:
	@echo -e '\033[1mCompiling CGO Qt Binding...\033[0m'
	cd standalone/qt/build && qmake ../qtbinding/qtbinding.pro && $(MAKE)
	cp -P $(QTBINDING_DLL) build/

qt-clean:
	cd standalone/qt/build && $(MAKE) clean
	rm standalone/qt/build/Makefile

crash-report:
	@echo -e '\033[1mCompiling Tool: Crash Report...\033[0m'
	cd support/tools/crash_report/build && qmake ../crash_report.pro && $(MAKE)
	cp support/tools/crash_report/build/$(QMAKE_RELEASE)/crash_report* build/

naive-debugger:
	@echo -e '\033[1mCompiling Tool: Naive Debugger...\033[0m'
	cd support/tools/naive_debugger/build && qmake ../naive_debugger.pro && $(MAKE)
	cp support/tools/naive_debugger/build/$(QMAKE_RELEASE)/naive_debugger* build/

tools: crash-report naive-debugger
	$(NOOP)

deps: check qt tools
	$(NOOP)

interpreter: deps
	@echo -e '\033[1mCompiling the Interpreter...\033[0m'
	go build $(ASAN_FLAG) -o build/$(EXENAME) main.go
	go build $(ASAN_FLAG) -buildmode=c-shared -o build/$(DLLNAME) main.lib.go
	chmod a+x build/$(DLLNAME)

release: qt-clean interpreter
ifeq ($(ASAN_ENABLED),1)
	@echo "*** cannot make release with ASAN enabled"
	exit 1
endif
	cd build && tar -vcapf release.tar.gz `ls * | grep -v -e "go_test_\|release"`


