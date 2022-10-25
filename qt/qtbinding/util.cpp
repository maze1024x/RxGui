#include <QApplication>
#include <QScreen>
#include "util.hpp"


QtString WrapString(QString str) {
    QString* ptr = new QString;
    *ptr = str;
    return { (void*) ptr };
}
QString UnwrapString(QtString str) {
    return QString(*(QString*)(str.ptr));
}
QIcon UnwrapIcon(QtIcon icon) {
    return QIcon(*(QIcon*)(icon.ptr));
}
QVariant UnwrapVariant(QtVariant v) {
    return QVariant(*(QVariant*)(v.ptr));
}

QStringList GetStringList(QtString* strings_ptr, size_t strings_len) {
    QStringList list;
    for (size_t i = 0; i < strings_len; i += 1) {
        list.append(UnwrapString(strings_ptr[i]));
    }
    return list;
}
QWidgetList GetWidgetList(void** widgets_ptr, size_t widgets_len) {
    QWidgetList list;
    for (size_t i = 0; i < widgets_len; i += 1) {
        list.append((QWidget*) widgets_ptr[i]);
    }
    return list;
}

QString EncodeBase64(QString str) {
    return QString::fromUtf8(str.toUtf8().toBase64(QByteArray::Base64UrlEncoding));
}
QString DecodeBase64(QString str) {
    return QString::fromUtf8(QByteArray::fromBase64(str.toUtf8(), QByteArray::Base64UrlEncoding));
}

int Get1remPixels() {
    static int value = -1;
    if (value > 0) {
        return value;
    }
    QScreen *screen = QGuiApplication::primaryScreen();
    QRect screenGeometry = screen->geometry();
    int screenHeight = screenGeometry.height();
    int screenWidth = screenGeometry.width();
    int minEdgeLength = std::min(screenHeight, screenWidth);
    value = ((RefScreen1remSize * minEdgeLength) / RefScreenMinEdgeLength);
    return value;
}
int GetScaledLength(int l) {
    if (l < 0) { 
        return 0;
    }
    return ((l * Get1remPixels()) / RefScreen1remSize);
}
QSize GetSizeFromRelative(QSize size_rem) {
    int unit = Get1remPixels();
    return QSize((unit * size_rem.width()), (unit * size_rem.height()));
}
void MoveToScreenCenter(QWidget* widget) {
    QScreen* screen = QGuiApplication::primaryScreen();
    widget->move(widget->pos() + (screen->geometry().center() - widget->geometry().center()));
}

QWidget* ObtainFocus() {
    return QApplication::focusWidget();
}
void RestoreFoucs(QWidget* w) {
    if (w != nullptr) {
        if (w != QApplication::focusWidget()) {
            w->setFocus(Qt::OtherFocusReason);
        }
    }
}

QMetaObject::Connection QtDynamicConnect (
        QObject* emitter , const QString& signalName,
        QObject* receiver, const QString& slotName
) {
    /* ref: https://stackoverflow.com/questions/26208851/qt-connecting-signals-and-slots-from-text */
    int index = emitter->metaObject()->indexOfSignal(QMetaObject::normalizedSignature(qPrintable(signalName)));
    if (index == -1) { return QMetaObject::Connection(); }
    QMetaMethod signal = emitter->metaObject()->method(index);
    index = receiver->metaObject()->indexOfSlot(QMetaObject::normalizedSignature(qPrintable(slotName)));
    if (index == -1) { return QMetaObject::Connection(); }
    QMetaMethod slot = receiver->metaObject()->method(index);
    return QObject::connect(emitter, signal, receiver, slot);
}

tdl_draw_t tdl_draw = nullptr;
void* workaround_tdl_vptr = nullptr;
void WorkaroundTdlDraw(QAbstractTextDocumentLayout* self, QPainter* painter, QAbstractTextDocumentLayout::PaintContext* context) {
    QTextDocument* doc = self->document();
    QTextOption opt = doc->defaultTextOption();
    opt.setWrapMode(QTextOption::WrapAtWordBoundaryOrAnywhere);
    self->blockSignals(true);
    doc->setDefaultTextOption(opt);
    self->blockSignals(false);
    tdl_draw(self, painter, context);
}


