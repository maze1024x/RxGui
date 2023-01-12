CONFIG += c++11
HEADERS += qtbinding.h util.hpp
SOURCES += api.cpp util.cpp
TARGET = qtbinding
TEMPLATE = lib
QT += widgets

RESOURCES += qtbinding.qrc Tango/

win32 {
    DEFINES += QTBINDING_WIN32_DLL
}

qtbinding_asan = $$(QTBINDING_ASAN)
equals(qtbinding_asan, "1") {
    message(address sanitizer enabled)
    CONFIG += sanitizer sanitize_address
}


