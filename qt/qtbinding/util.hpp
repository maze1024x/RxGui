#ifndef UTIL_HPP
#define UTIL_HPP

#include <QCoreApplication>
#include <QObject>
#include <QMetaMethod>
#include <QWidget>
#include <QLayout>
#include <QSplitter>
#include <QEvent>
#include <QPointer>
#include <QMenu>
#include <QDialog>
#include <QDialogButtonBox>
#include <QComboBox>
#include <QLabel>
#include <QTextDocument>
#include <QAbstractTextDocumentLayout>
#include <QFrame>
#include <QPainter>
#include <QTextLayout>
#include <QTextLine>
#include <QTextEdit>
#include <QPlainTextEdit>
#include <QLineEdit>
#include <QToolTip>
#include <QHeaderView>
#include <QTreeWidget>
#include <QListWidget>
#include <QScrollBar>
#include <QScrollArea>
#include <QTimer>
#include <QDebug>
#include <cstdlib>
#include <cmath>
#include "qtbinding.h"


#define RefScreen1remSize 16
#define RefScreenMinEdgeLength 768

QtString WrapString(QString str);
QString UnwrapString(QtString str);
QIcon UnwrapIcon(QtIcon icon);
QVariant UnwrapVariant(QtVariant v);

QStringList GetStringList(QtString* strings_ptr, size_t strings_len);
QWidgetList GetWidgetList(void** widgets_ptr, size_t widgets_len);

QString EncodeBase64(QString str);
QString DecodeBase64(QString str);

int Get1remPixels();
int GetScaledLength(int l);
QSize GetSizeFromRelative(QSize size_rem);
void MoveToScreenCenter(QWidget* widget);

QWidget* ObtainFocus();
void RestoreFoucs(QWidget* w);


typedef void (*callback_t)(uint64_t);
Q_DECLARE_METATYPE(callback_t);

class CallbackObject;

QMetaObject::Connection QtDynamicConnect (
        QObject* emitter , const QString& signalName,
        QObject* receiver, const QString& slotName
);

class CallbackExecutor final: public QObject {
    Q_OBJECT
public:
    CallbackExecutor(QWidget *parent = nullptr): QObject(parent) {
        QObject::connect (
            this, &CallbackExecutor::QueueCallback,
            this, &CallbackExecutor::__InvokeCallback,
            Qt::QueuedConnection
        );
    };
    virtual ~CallbackExecutor() {};
signals:
    void QueueCallback(callback_t cb, uint64_t payload);
protected slots:
    void __InvokeCallback(callback_t cb, uint64_t payload) {
        cb(payload);
    };
};

class CallbackObject final: public QObject {
    Q_OBJECT
public:
    callback_t cb;
    uint64_t payload;
    CallbackObject(QObject* parent, callback_t cb, uint64_t payload): QObject(parent) {
        this->cb = cb;
        this->payload = payload;
    };
    virtual ~CallbackObject() {};
public slots:
    void Slot() {
        cb(payload);
    };
};

class EventListener final: public QObject {
    Q_OBJECT
public:
    QEvent::Type accept_type;
    bool prevent_default;
    callback_t cb;
    uint64_t payload;
    QObject* current_object;
    QEvent* current_event;
    EventListener(QEvent::Type t, bool prevent, callback_t cb, uint64_t payload): QObject(nullptr) {
        this->accept_type = t;
        this->prevent_default = prevent;
        this->cb = cb;
        this->payload = payload;
    }
    bool eventFilter(QObject* obj, QEvent *event) override {
        if (event->type() == accept_type) {
            current_object = obj;
            current_event = event;
            cb(payload);
            current_event = nullptr;
            current_object = nullptr;
            if (prevent_default) {
                event->ignore();
                return true; // stops propagation
            } else {
                return false;
            }
        } else {
            return false;
        }
    }
};

class MenuHandle: public QObject {
    Q_OBJECT
public:
    MenuHandle(QMenu* menu, QObject* parent): QObject(parent), menu(menu) {}
    virtual ~MenuHandle() {
        delete menu;
    }
protected:
    QMenu* menu;
};

class ContextMenuBinding: public QObject {
    Q_OBJECT
protected:
    ContextMenuBinding(QWidget* widget, QMenu* menu, Qt::ContextMenuPolicy policy): QObject(nullptr), widget(widget), menu(menu), policy(policy) {}
public:
    static ContextMenuBinding* TryCreate(QWidget* widget, QMenu* menu) {
        Qt::ContextMenuPolicy policy = widget->contextMenuPolicy();
        if (policy == Qt::CustomContextMenu) {
            return nullptr;
        }
        widget->setContextMenuPolicy(Qt::CustomContextMenu);
        ContextMenuBinding* b = new ContextMenuBinding(widget, menu, policy);
        QObject::connect(widget, &QWidget::customContextMenuRequested, b, &ContextMenuBinding::Popup);
        return b;
    }
    virtual ~ContextMenuBinding() {
        if (widget) {
            widget->setContextMenuPolicy(policy);
        }
        if (menu) {
            menu->hide();
        }
    }
public slots:
    void Popup(const QPoint& p) {
        if (widget && menu) {
            QWidget* w = nullptr;
            QAbstractScrollArea* a = qobject_cast<QAbstractScrollArea*>(widget);
            if (a != nullptr) {
                w = a->viewport();
            } else {
                w = widget;
            }
            QPoint g = w->mapToGlobal(p);
            menu->popup(g);
        }
    }
protected:
    QPointer<QWidget> widget;
    QPointer<QMenu> menu;
    Qt::ContextMenuPolicy policy;
};

class SizeHintWidget: public QWidget {
    Q_OBJECT
public:
    SizeHintWidget(QSize s, QWidget* parent): QWidget(parent), specifiedSizeHint(s) {}
    virtual ~SizeHintWidget() {}
    QSize sizeHint() const override {
        return specifiedSizeHint;
    }
protected:
    QSize specifiedSizeHint;
};

class WeakWidget: public QWidget {
    Q_OBJECT
public:
    WeakWidget(QLayout* layout): QWidget(nullptr) {
        setLayout(layout);
    }
    virtual ~WeakWidget() {
        QObjectList dup;
        for (QObject* child: children()) {
            dup.append(child);
        }
        for (QObject* child: dup) {
            QWidget* widget = qobject_cast<QWidget*>(child);
            if (widget != nullptr) {
                widget->hide();
                widget->setParent(nullptr);
            }
        }
    }
    static WeakWidget* Wrap(QWidget* w, bool check = true) {
        QLayout* layout = new QHBoxLayout();
        if (!(check) || w->parentWidget() == nullptr) {
            layout->addWidget(w);
        }
        layout->setContentsMargins(0, 0, 0, 0);
        WeakWidget* wrapper = new WeakWidget(layout);
        return wrapper;
    }
    static QWidget* Inner(QWidget* w) {
        if (w != nullptr) {
            WeakWidget* wrapper = qobject_cast<WeakWidget*>(w);
            if (wrapper != nullptr) {
                QLayoutItem* item = wrapper->layout()->itemAt(0);
                if (item != nullptr) {
                    QWidget* inner = item->widget();
                    if (inner != nullptr) {
                        return inner;
                    }
                }
            }
        }
        return nullptr;
    }
    void paintEvent(QPaintEvent* ev) override {
        QStyleOption opt;
        opt.initFrom(this);
        QPainter p(this);
        style()->drawPrimitive(QStyle::PE_Widget, &opt, &p, this);
        QWidget::paintEvent(ev);
    }
};

class DynamicWidget: public QWidget {
    Q_OBJECT
protected:
    QWidget* current;
public:
    DynamicWidget(QWidget* parent): QWidget(parent) {
        current = nullptr;
        QLayout* layout = new QHBoxLayout();
        layout->setContentsMargins(0, 0, 0, 0);
        setLayout(layout);
    }
    virtual ~DynamicWidget() {
        setWidget(nullptr);
    }
    void setWidget(QWidget* w) {
        if (w == current) {
            return;
        }
        if (w != nullptr && w->parentWidget() != nullptr) {
            return;
        }
        if (current != nullptr) {
            current->hide();
            current->setParent(nullptr);
            layout()->removeWidget(current);
        }
        if (w != nullptr) {
            layout()->addWidget(w);
            w->show();
            connect(w, &QObject::destroyed, this, &DynamicWidget::widgetDestroyed, Qt::UniqueConnection);
        }
        current = w;
    }
public slots:
    void widgetDestroyed(QObject* w) {
        if (static_cast<QObject*>(current) == w) {
            layout()->removeWidget(current);
            current = nullptr;
        }
    }
};

class ComboBoxDialog: public QDialog {
    Q_OBJECT
public:
    ComboBoxDialog(QString prompt, QWidget* parent): QDialog(parent), comboBox(new QComboBox()) {
        // reference: https://code.woboq.org/qt5/qtbase/src/widgets/dialogs/qinputdialog.cpp.html#_ZN19QInputDialogPrivate12ensureLayoutEv
        QLabel* label = new QLabel(prompt);
        label->setBuddy(comboBox);
        label->setSizePolicy(QSizePolicy::Minimum, QSizePolicy::Fixed);
        QDialogButtonBox* buttonBox = new QDialogButtonBox(QDialogButtonBox::Ok | QDialogButtonBox::Cancel, Qt::Horizontal);
        QObject::connect(buttonBox, &QDialogButtonBox::accepted, this, &QDialog::accept);
        QObject::connect(buttonBox, &QDialogButtonBox::rejected, this, &QDialog::reject);
        QLayout* layout = new QVBoxLayout();
        layout->setSizeConstraint(QLayout::SetMinAndMaxSize);
        layout->addWidget(label);
        layout->addWidget(comboBox);
        layout->addWidget(buttonBox);
        setLayout(layout);
    }
    QComboBox* ComboBox() {
        return comboBox;
    }
protected:
    QComboBox* comboBox;
};

class CustomDialog: public QDialog {
    Q_OBJECT
public:
    CustomDialog(): QDialog(nullptr) {
        setWindowFlags(windowFlags() & ~Qt::WindowContextHelpButtonHint);
    }
    virtual ~CustomDialog() {}
    virtual void closeEvent(QCloseEvent* ev) override {
        QWidget::closeEvent(ev);
    }
    virtual void keyPressEvent(QKeyEvent* ev) override {
        if (ev->matches(QKeySequence::Cancel)) {
            QCoreApplication::postEvent(this, new QCloseEvent());
        } else {
            QWidget::keyPressEvent(ev);
        }
    }
public slots:
    virtual void done(int) override {
        // no-op
    }
};

class SmartScrollArea final: public QScrollArea {
    Q_OBJECT
public:
    enum Direction { BothDirection, VerticalOnly, HorizontalOnly };
    SmartScrollArea(Direction d, QWidget* parent): QScrollArea(parent), direction(d) {
        setWidgetResizable(true);
        if (direction == VerticalOnly) {
            setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
        } else if (direction == HorizontalOnly) {
            setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
        }
    }
    virtual ~SmartScrollArea() {}
protected:
    Direction direction;
    void resizeEvent(QResizeEvent* event) override {
        if (direction == VerticalOnly) {
            widget()->setMaximumWidth(event->size().width());
        } else if (direction == HorizontalOnly) {
            widget()->setMaximumHeight(event->size().height());
        }
        QScrollArea::resizeEvent(event);
    }
public:
    bool TextViewShouldHaveMinWidthHint() {
        return (direction != VerticalOnly);
    }
};

// modified from: https://doc.qt.io/qt-5/qtwidgets-widgets-elidedlabel-example.html
class ElidedLabel: public QFrame {
    Q_OBJECT
    Q_PROPERTY(QString text READ text WRITE setText)
    Q_PROPERTY(bool isElided READ isElided)
protected:
    bool elided;
    QString content;
    bool multiline;
    struct Style {
        Style(): bold(false), italic(false), underline(false), strikeOut(false), color("") {}
        bool bold; bool italic; bool underline; bool strikeOut; QString color;
    };
    Style style;
public:
    ElidedLabel(QString text, QWidget* parent)
        : QFrame(parent)
        , elided(false)
        , content(text)
        , multiline(false)
        , style(Style())
    {
        setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Preferred);
    }
    virtual ~ElidedLabel() {}
    void setText(const QString &newText) {
        content = newText;
        update();
    }
    const QString & text() const { return content; }
    bool isElided() const { return elided; }
    virtual QSize sizeHint() const override {
        return fontMetrics().size(Qt::TextSingleLine, content);
    }
    virtual QSize minimumSizeHint() const override {
        return fontMetrics().size(Qt::TextSingleLine, "...");
    }
protected:
    void updateStyle() {
        // ad hoc implementation
        style = Style();
        QString qss = styleSheet();
        qss = qss.replace("{", ";");
        qss = qss.replace("}", ";");
        QStringList segments = qss.split(";");
        for (QString s: segments) {
            if (s.contains("font-weight")) {
                style.bold = s.contains("bold");
            }
            if (s.contains("font-style")) {
                style.italic = s.contains("italic");
            }
            if (s.contains("text-decoration")) {
                style.underline = s.contains("underline");
                style.strikeOut = s.contains("line-through");
            }
            if (s.contains("color")) {
                QStringList kv = s.split(":");
                if (1 < kv.length()) {
                    style.color = kv[1];
                }
            }
        }
    }
    void paintEvent(QPaintEvent *event) override {
        QFrame::paintEvent(event);
        QPainter painter(this);
        updateStyle(); {
            QFont font = painter.font();
            font.setBold(style.bold);
            font.setItalic(style.italic);
            font.setUnderline(style.underline);
            font.setStrikeOut(style.strikeOut);
            painter.setFont(font);
            if (style.color != "") {
                painter.setPen(QColor(style.color));
            }
        }
        QFontMetrics fontMetrics = painter.fontMetrics();
        bool didElide = false;
        int lineSpacing = fontMetrics.lineSpacing();
        if (multiline) {
            int y = 0;
            QTextLayout textLayout(content, painter.font());
            textLayout.beginLayout();
            forever {
                QTextLine line = textLayout.createLine();
                if (!line.isValid())
                    break;
                line.setLineWidth(width());
                int nextLineY = y + lineSpacing;
                if (height() >= nextLineY + lineSpacing) {
                    line.draw(&painter, QPoint(0, y));
                    y = nextLineY;
                } else {
                    QString lastLine = content.mid(line.textStart());
                    QString elidedLastLine = fontMetrics.elidedText(lastLine, Qt::ElideRight, width());
                    painter.drawText(QPoint(0, y + fontMetrics.ascent()), elidedLastLine);
                    didElide = (elidedLastLine != lastLine);
                    break;
                }
            }
            textLayout.endLayout();
        } else {
            int y = ((height() - lineSpacing) / 2);
            QString elidedContent = fontMetrics.elidedText(content, Qt::ElideRight, width());
            painter.drawText(QPoint(0, y + fontMetrics.ascent()), elidedContent);
            didElide = (elidedContent != content);
        }
        if (didElide != elided) {
            elided = didElide;
            if (elided) {
                setToolTip(content);
            } else {
                setToolTip("");
            }
            emit elisionChanged(didElide);
        }
    }
signals:
    void elisionChanged(bool elided);
};

class TextView: public QLabel {
    Q_OBJECT
public:
    TextView(QString text, Qt::TextFormat format, QWidget* parent): QLabel(parent) {
        setWordWrap(true);
        setTextFormat(format);
        setText(text);
        setAlignment(Qt::AlignLeft | Qt::AlignTop);
        setTextInteractionFlags(Qt::TextSelectableByMouse | textInteractionFlags());
    }
    virtual ~TextView() {}
    virtual QSize minimumSizeHint() const override {
        bool scroll_area_workaround = false;
        QWidget* w = parentWidget();
        while (w != nullptr && qobject_cast<WeakWidget*>(w) != nullptr) {
            w = w->parentWidget();
        }
        if (w != nullptr) {
            if (w->objectName() == "qt_scrollarea_viewport") {
                w = w->parentWidget();
            }
            if (SmartScrollArea* a = qobject_cast<SmartScrollArea*>(w)) {
                if (a->TextViewShouldHaveMinWidthHint()) {
                    scroll_area_workaround = true;
                }
            }
        }
        if (scroll_area_workaround) {
            return QLabel::minimumSizeHint();
        }
        // prevent long word from contributing a minimum width
        QSize size = QLabel::minimumSizeHint();
        size.setWidth(0);
        return size;
    }
};

// modified from: https://stackoverflow.com/questions/27000484/add-custom-widgets-as-qtablewidget-horizontalheader
class CustomHeaderView: public QHeaderView {
    Q_OBJECT
protected:
    QMap<int,QWidget*> sectionsWidgets;
public:
    CustomHeaderView(QWidget* parent): QHeaderView(Qt::Horizontal, parent) {
        connect(this, &QHeaderView::geometriesChanged, this, &CustomHeaderView::handleGeometriesChanged);
        connect(this, &QHeaderView::sectionResized,    this, &CustomHeaderView::handleSectionResized);
        connect(this, &QHeaderView::sectionMoved,      this, &CustomHeaderView::handleSectionMoved);
    }
    ~CustomHeaderView() {}
    void setSectionWidget(int index, QWidget * widget) {
        if (sectionsWidgets[index] != nullptr) {
            delete sectionsWidgets[index];
        }
        widget->setParent(this);
        sectionsWidgets[index] = widget;
    }
    void workaround() {
        placeAllWidgets();
    }
protected:
    void showEvent(QShowEvent* e) override {
        for (int i = 0; i < count(); i += 1) {
            if (sectionsWidgets[i] == nullptr) {
                sectionsWidgets[i] = new QWidget(this);
            }
            placeWidget(i);
            sectionsWidgets[i]->show();
        }
        QHeaderView::showEvent(e);
    }
    QSize sectionSizeFromContents(int logical) const override {
        return sectionsWidgets[logical]->sizeHint();
    }
    void placeWidget(int logical) {
        sectionsWidgets[logical]->setGeometry(
            sectionViewportPosition(logical), 0,
            sectionSize(logical), height()
        );
    }
    void placeAllWidgets() {
        for (int j = 0; j < count(); j += 1) {
            placeWidget(logicalIndex(j));
        }
    }
protected slots:
    void handleGeometriesChanged() {
        placeAllWidgets();
    }
    void handleSectionResized(int i) {
        for (int j = visualIndex(i); j < count(); j += 1) {
            placeWidget(logicalIndex(j));
        }
    }
    void handleSectionMoved(int logical, int oldVisual, int newVisual) {
        Q_UNUSED(logical);
        for (int j = qMin(oldVisual, newVisual); j < count(); j += 1) {
            placeWidget(logicalIndex(j));
        }
    }
public slots:
    void handleScrollBarValueChanged() {
        placeAllWidgets();
    }
};

class ListWidgetInterface {
public:
    enum InsertHint { Prepend, Append, InsertAbove, InsertBelow };
    enum MoveHint { Up, Down };
    virtual QWidget* CastToWidget() = 0;
    virtual QString Current(bool* exists) = 0;
    virtual QString Activation() = 0;
    virtual QStringList All() = 0; // returned keys must be in list order
    virtual QStringList Selection() = 0; // (same as the comment above)
    virtual QStringList ContiguousSelection() = 0; // (same as the comment above)
    virtual void Insert(InsertHint h, QString pivot, QString key, QWidgetList widgets) = 0;
    virtual void Update(QString key) = 0;
    virtual bool Move(MoveHint h, QString key) = 0;
    virtual void Delete(QString key) = 0;
    virtual void Reorder(QStringList order) = 0;
};

class DefaultListWidget: public QTreeWidget, public ListWidgetInterface {
    Q_OBJECT
public:
    using Select = QAbstractItemView::SelectionMode;
    DefaultListWidget(size_t columns, Select select, QWidgetList headers, int stretch, QWidget* parent): QTreeWidget(parent) {
        if (columns == 0) {
            columns = 1;
            setHeaderHidden(true);
        }
        setColumnCount(columns);
        if (select == QAbstractItemView::NoSelection) {
            setFocusPolicy(Qt::NoFocus);
        }
        setSelectionMode(select);
        setRootIsDecorated(false);
        QObject::connect(this, &QTreeWidget::itemSelectionChanged, this, [this] () -> void {
            if (selectedItems().isEmpty()) {
                setCurrentItem(nullptr);
            }
        });
        if (headers.length() > 0) {
            CustomHeaderView* header_view = new CustomHeaderView(this);
            int l = qMin(columnCount(), headers.length());
            for (int i = 0; i < l; i += 1) {
                headerItem()->setData(i, Qt::DisplayRole, "");
                header_view->setSectionWidget(i, WeakWidget::Wrap(headers[i]));
            }
            setHeader(header_view);
            connect(horizontalScrollBar(), &QScrollBar::valueChanged, header_view, &CustomHeaderView::handleScrollBarValueChanged);
        }
        header()->setSectionsMovable(false);
        header()->setStretchLastSection(false);
        if (stretch < 0) {
            header()->setSectionResizeMode(QHeaderView::Stretch);
        } else {
            if (stretch < columnCount()) {
                header()->setSectionResizeMode(QHeaderView::ResizeToContents);
                header()->setSectionResizeMode(stretch, QHeaderView::Stretch);
            }
        }
        connect(this, &QAbstractItemView::doubleClicked, this, &DefaultListWidget::DoubleClickActivate);
    }
    virtual ~DefaultListWidget() {}
protected:
    QString activation;
    void DoubleClickActivate(const QModelIndex& index) {
        activation = GetKey(itemFromIndex(index));
        emit activationTriggered();
    }
    void ReturnPressActivate() {
        bool exists;
        QString current = Current(&exists);
        if (exists) {
            activation = current;
            emit activationTriggered();
        }
    }
    virtual void keyPressEvent(QKeyEvent* ev) override {
        if (ev->key() == Qt::Key_Return) {
            ReturnPressActivate();
        }
        QTreeWidget::keyPressEvent(ev);
    }
    virtual void focusInEvent(QFocusEvent* ev) override {
        // bypass QAbstractItemView::focusInEvent
        QAbstractScrollArea::focusInEvent(ev);
    }
    virtual void focusOutEvent(QFocusEvent* ev) override {
        // bypass QAbstractItemView::focusOutEvent
        QAbstractScrollArea::focusOutEvent(ev);
    }
    int Size() {
        return topLevelItemCount();
    }
    QString At(int index) {
        return GetKey(topLevelItem(index));
    }
    static QString GetKey(QTreeWidgetItem* item) {
        return item->data(0, Qt::UserRole).toString();
    }
    static void SetKey(QTreeWidgetItem* item, QString key) {
        item->setData(0, Qt::UserRole, key);
    }
    void SetWidgets(QTreeWidgetItem* item, QWidgetList widgets, bool check = true) {
        int l = qMin(columnCount(), widgets.length());
        for (int i = 0; i < l; i += 1) {
            setItemWidget(item, i, WeakWidget::Wrap(widgets[i], check));
        }
    }
    QWidgetList ObtainWidgets(QTreeWidgetItem* item) {
        QWidgetList widgets;
        for (int i = 0; i < columnCount(); i += 1) {
            QWidget* inner = WeakWidget::Inner(itemWidget(item, i));
            if (inner != nullptr) {
                widgets.append(inner);
            } else {
                break;
            }
        }
        return widgets;
    }
    bool LookupItemPos(QString key, int* out) {
        for (int i = 0; i < Size(); i += 1) {
            if (At(i) == key) {
                *out = i;
                return true;
            }
        }
        return false;
    }
    void UpdateColumnsWidth() {
        for (int i = 0; i < columnCount(); i += 1) {
            QHeaderView::ResizeMode mode = header()->sectionResizeMode(i);
            if (mode != QHeaderView::Stretch) {
                resizeColumnToContents(i);
            }
        }
        CustomHeaderView* chv = qobject_cast<CustomHeaderView*>(header());
        if (chv != nullptr) {
            chv->workaround();
        }
    }
signals:
    void activationTriggered();
public:
    virtual QWidget* CastToWidget() override {
        return static_cast<QWidget*>(this);
    }
    virtual QString Current(bool* exists) override {
        QTreeWidgetItem* item = currentItem();
        if (item != nullptr) {
            *exists = true;
            return (GetKey(item));
        } else {
            *exists = false;
            return "";
        }
    }
    virtual QString Activation() override {
        return activation;
    }
    virtual QStringList All() override {
        QStringList keys;
        for (int i = 0; i < Size(); i += 1) {
            keys.append(At(i));
        }
        return keys;
    }
    virtual QStringList Selection() override {
        QStringList keys;
        for (int i = 0; i < Size(); i += 1) {
            QTreeWidgetItem* item = topLevelItem(i);
            if (item->isSelected()) {
                keys.append(GetKey(item));
            }
        }
        return keys;
    }
    virtual QStringList ContiguousSelection() override {
        int first = -1;
        int last = -1;
        bool has_gap = false;
        for (int i = 0; i < Size(); i += 1) {
            QTreeWidgetItem* item = topLevelItem(i);
            if (item->isSelected()) {
                if (first < 0) {
                    first = i;
                    last = i;
                } else {
                    if (has_gap) {
                        return QStringList();
                    } else {
                        last = i;
                    }
                }
            } else {
                if (first >= 0) {
                    has_gap = true;
                }
            }
        }
        QStringList keys;
        if (first >= 0) {
            for (int i = first; i <= last; i += 1) {
                keys.append(At(i));
            }
        }
        return keys;
    }
    virtual void Insert(InsertHint h, QString pivot, QString key, QWidgetList widgets) override {
        int pos = 0;
        if (h == Prepend) {
            pos = 0;
        } else if (h == Append) {
            pos = Size();
        } else if (h == InsertAbove) {
            int i = 0;
            bool ok = LookupItemPos(pivot, &i);
            if (!(ok)) { return; }
            pos = i;
        } else if (h == InsertBelow) {
            int i = 0;
            bool ok = LookupItemPos(pivot, &i);
            if (!(ok)) { return; }
            pos = (i + 1);
        }
        Q_ASSERT(pos >= 0);
        QTreeWidgetItem* item = new QTreeWidgetItem();
        SetKey(item, key);
        insertTopLevelItem(pos, item); // note: we assume the key is unique
        SetWidgets(item, widgets);
        UpdateColumnsWidth();
    }
    virtual void Update(QString key) override {
        static_cast<void>(key);
        QTimer::singleShot(0, this, [this] () -> void {
            UpdateColumnsWidth();
        });
    }
    virtual bool Move(MoveHint h, QString key) override {
        int from = 0;
        bool ok = LookupItemPos(key, &from);
        if (!(ok)) { return false; }
        int to = 0;
        if (h == Up) {
            to = (from - 1);
        } else if (h == Down) {
            to = (from + 1);
        }
        return MoveByIndex(from, to);
    }
    virtual void Delete(QString key) override {
        int i;
        bool ok = LookupItemPos(key, &i);
        if (!(ok)) { return; }
        QTreeWidgetItem* item = topLevelItem(i);
        for (int j = 0; j < columnCount(); j += 1) {
            removeItemWidget(item, j);
        }
        if (item == currentItem()) {
            setCurrentItem(nullptr);
        }
        delete item;
        UpdateColumnsWidth();
    }
    virtual void Reorder(QStringList order) override {
        QMap<QString,int> mapping;
        for (int i = 0; i < Size(); i += 1) {
            QTreeWidgetItem* item = topLevelItem(i);
            QString key = item->data(0, Qt::UserRole).toString();
            mapping.insert(key, i);
        }
        for (int i = 0; i < order.length(); i += 1) {
            QString key = order[i];
            auto it = mapping.find(key);
            if (it != mapping.end()) {
                int index = it.value();
                if (index != i) {
                    QTreeWidgetItem* i_item = topLevelItem(i);
                    QString i_key = i_item->data(0, Qt::UserRole).toString();
                    MoveByIndex(index, i);
                    if ((i+1) != index) {
                        MoveByIndex((i+1), index);
                    }
                    mapping[i_key] = index;
                }
                mapping.remove(key);
            }
        }
    }
    bool MoveByIndex(int from, int to) {
        /* reference: https://stackoverflow.com/questions/2035932/raise-and-lower-qtreewidgetitem-in-a-qtreewidget */
        /* reference: https://stackoverflow.com/questions/25559221/qtreewidgetitem-issue-items-set-using-setwidgetitem-are-dispearring-after-movin */
        QTreeWidgetItem* item = topLevelItem(from);
        if (0 <= to && to < Size()) {
            auto focus = ObtainFocus();
            auto widgets = ObtainWidgets(item);
            auto current = (item == currentItem());
            auto selected = item->isSelected();
            item->setSelected(false); // prevents strange selection afterimage
            auto data = model()->itemData(model()->index(from, 0));
            model()->removeRow(from);
            model()->insertRow(to);
            model()->setItemData(model()->index(to, 0), data);
            item = topLevelItem(to);
            item->setSelected(selected);
            if (current) { setCurrentItem(item, 0, QItemSelectionModel::NoUpdate); }
            SetWidgets(item, widgets, false);
            RestoreFoucs(focus);
            return true;
        } else {
            return false;
        }
    }
};

#endif  // UTIL_HPP


