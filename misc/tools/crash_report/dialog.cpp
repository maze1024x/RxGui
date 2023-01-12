#include <QTextStream>
#include "dialog.h"
#include "ui_dialog.h"

Dialog::Dialog(QString program_path, QWidget *parent)
    : QDialog(parent)
    , ui(new Ui::Dialog)
{
    ui->setupUi(this);
    setWindowFlag(Qt::WindowStaysOnTopHint);
    QString info;
    QTextStream in(stdin);
    while(!in.atEnd()) {
        QString line = in.readLine();
        info += line;
        info += "\n";
    }
    info += "*** END ***";
    ui->contentText->setHtml(QString(
        "<b>Program:</b>"
        "<pre>%1</pre>"
        "<b>Info:</b>"
        "%2"
    ).arg(program_path.toHtmlEscaped(), info));
}

Dialog::~Dialog()
{
    delete ui;
}

