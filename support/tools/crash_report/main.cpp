#include "dialog.h"

#include <QApplication>

int main(int argc, char *argv[])
{
    QApplication a(argc, argv);
    QString program_path;
    if (argc >= 2) {
        program_path = QString(argv[1]);
    }
    Dialog w(program_path);
    w.show();
    return a.exec();
}
