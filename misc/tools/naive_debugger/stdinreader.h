#ifndef STDINREADER_H
#define STDINREADER_H

#include <QThread>
#include <QTextStream>

class StdinReader: public QThread {
    Q_OBJECT
public:
    StdinReader(): QThread() {}
    virtual ~StdinReader() {}
    void run() override {
        QTextStream in(stdin);
        while(!in.atEnd()) {
            QString line = in.readLine();
            emit NewLine(line);
        }
    }
signals:
    void NewLine(QString line);
};

#endif // STDINREADER_H
