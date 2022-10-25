#include <QScrollBar>
#include <QKeyEvent>
#include "mainwindow.h"
#include "ui_mainwindow.h"

void TextAppend(QTextBrowser* t, QString content) {
    // reference: stackoverflow.com/questions/54501745/performantly-appending-rich-text-into-qtextedit-or-qtextbrowser-in-qt/54501760
    QScrollBar* bar = t->verticalScrollBar();
    const bool should_scroll = (
        (!(bar->isVisible()) && t->isVisible())
        || (bar->value() == bar->maximum())
    );
    QTextDocument* doc = t->document();
    QTextCursor cursor(doc);
    cursor.movePosition(QTextCursor::End);
    cursor.insertHtml(content + "<br>");
    if (should_scroll) {
        bar->setValue(bar->maximum());
    }
}

void ZoomIn(QTextBrowser* t) {
    for (int i = 0; i < 4; i += 1) {
        t->zoomIn();
    }
}

MainWindow::MainWindow(QString program_path, QWidget *parent)
    : QMainWindow(parent)
    , history(QVector<QString>(1, "")), ptr(0)
    , ui(new Ui::MainWindow)
{
    ui->setupUi(this);
    if (program_path != "") {
        setWindowTitle(windowTitle() + " - " + program_path);
    }
    connect(ui->tabWidget, &QTabWidget::currentChanged, this, [this] (int index) -> void {
        QWidget* widget = ui->tabWidget->widget(index);
        if (widget == ui->everything || widget == ui->repl) {
            ui->cmdInput->setVisible(true);
        } else {
            ui->cmdInput->setVisible(false);
        }
    });
    connect(ui->cmdInput, &QLineEdit::returnPressed, this, [this] () -> void {
        QString cmd = ui->cmdInput->text();
        history[history.length()-1] = cmd;
        ptr = history.length();
        history.append("");
        ui->cmdInput->setText("");
        emit Command(cmd);
    });
    ui->cmdInput->installEventFilter(this);
    ui->cmdInput->setFocus(Qt::OtherFocusReason);
    ZoomIn(ui->everythingText);
    ZoomIn(ui->replText);
    ZoomIn(ui->inspectionText);
    ZoomIn(ui->ioText);
}

bool MainWindow::eventFilter(QObject* obj, QEvent* ev) {
    if (obj == ui->cmdInput) {
        if (ev->type() == QEvent::KeyPress) {
            int key = static_cast<QKeyEvent*>(ev)->key();
            if (key == Qt::Key_Up) {
                if (0 <= ptr-1 && ptr-1 < history.length()) {
                    if (ptr+1 == history.length()) {
                        history[ptr] = ui->cmdInput->text();
                    }
                    ptr--;
                    ui->cmdInput->setText(history[ptr]);
                }
                return true;
            } else if (key == Qt::Key_Down) {
                if (0 <= ptr+1 && ptr+1 < history.length()) {
                    ptr++;
                    ui->cmdInput->setText(history[ptr]);
                }
                return true;
            }
        }
    }
    return QMainWindow::eventFilter(obj, ev);
}

void MainWindow::Message(QString category, QString content) {
    if (category == "<control>") {
        if (content == "PROGRAM_CRASH") {
            ui->cmdInput->setDisabled(true);
            ui->cmdInput->setText("");
        }
    }
    if (category == "repl" || category == "*") {
        TextAppend(ui->replText, content);
    }
    if (category == "inspection" || category == "*") {
        TextAppend(ui->inspectionText, content);
    }
    if (category == "io" || category == "*") {
        TextAppend(ui->ioText, content);
    }
    if (category != "<control>") {
        TextAppend(ui->everythingText, content);
    }
}

MainWindow::~MainWindow()
{
    delete ui;
}

