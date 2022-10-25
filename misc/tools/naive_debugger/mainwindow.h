#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QVector>
#include <QEvent>

QT_BEGIN_NAMESPACE
namespace Ui { class MainWindow; }
QT_END_NAMESPACE

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    MainWindow(QString program_path, QWidget *parent = nullptr);
    ~MainWindow();

protected:
    QVector<QString> history;
    int ptr;
    bool eventFilter(QObject* obj, QEvent* ev) override;

signals:
    void Command(QString cmd);

public slots:
    void Message(QString category, QString content);

private:
    Ui::MainWindow *ui;
};
#endif // MAINWINDOW_H
