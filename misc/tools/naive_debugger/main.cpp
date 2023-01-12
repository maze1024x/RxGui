#include <QApplication>
#include <QThread>
#include <QTextStream>
#include <QSocketNotifier>
#include <QJsonDocument>
#ifdef _WIN32
#include <windows.h>
#endif
#include "mainwindow.h"
#include "stdinreader.h"

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    QString program_path;
    if (argc >= 2) {
        program_path = QString(argv[1]);
    }
    MainWindow w(program_path);
    StdinReader r;
    QObject::connect(&r, &StdinReader::NewLine, &w, [&] (QString line) -> void {
        QJsonDocument doc = QJsonDocument::fromJson(line.toUtf8());
        if (!(doc.isNull())) {
            QString category = doc["category"].toString();
            QString content = doc["content"].toString();
            if (category != "" && content != "") {
                w.Message(category, content);
            }
        }
    });
    r.start();
    QTextStream out(stdout);
    QObject::connect(&w, &MainWindow::Command, &w, [&] (QString cmd) -> void {
        out << cmd << "\n";
        out.flush();
    });
    w.show();
    #ifdef _WIN32
    // raise window
    HWND hwnd = (HWND) w.winId();
    SetWindowPos(hwnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOACTIVATE);
    SetWindowPos(hwnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOACTIVATE);
    #endif
    int code = a.exec();
    r.terminate();
    r.wait();
    return code;
}
