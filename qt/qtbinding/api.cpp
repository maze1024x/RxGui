#include <QApplication>
#include <QClipboard>
#include <QScreen>
#include <QDir>
#include <QMainWindow>
#include <QGroupBox>
#include <QInputDialog>
#include <QMessageBox>
#include <QFileDialog>
#include <QBuffer>
#include <QByteArray>
#include <QString>
#include <QVector>
#include <QVariant>
#include <QResizeEvent>
#include <QGridLayout>
#include <QAction>
#include <QActionGroup>
#include <QMenu>
#include <QMenuBar>
#include <QToolBar>
#include <QToolButton>
#include <QIcon>
#include <QPixmap>
#include <QLabel>
#include <QCheckBox>
#include <QLineEdit>
#include <QPlainTextEdit>
#include <QPushButton>
#include <QSlider>
#include <QProgressBar>
#include <QVariantMap>
#include <QVariantList>
#include <QTimer>
#include <QUuid>
#ifdef _WIN32
#include <windows.h>
#endif
#include "util.hpp"
#include "qtbinding.h"


const int QtEventMove = QEvent::Move;
const int QtEventResize = QEvent::Resize;
const int QtEventShow = QEvent::Show;
const int QtEventClose = QEvent::Close;
const int QtEventFocusIn = QEvent::FocusIn;
const int QtEventFocusOut = QEvent::FocusOut;
const int QtEventWindowActivate = QEvent::WindowActivate;
const int QtEventWindowDeactivate = QEvent::WindowDeactivate;
const int QtEventDynamicPropertyChange = QEvent::DynamicPropertyChange;

const int QtAlignDefault = Qt::Alignment();
const int QtAlignLeft = Qt::AlignLeft;
const int QtAlignRight = Qt::AlignRight;
const int QtAlignHCenter = Qt::AlignHCenter;
const int QtAlignTop = Qt::AlignTop;
const int QtAlignBottom = Qt::AlignBottom;
const int QtAlignVCenter = Qt::AlignVCenter;

const int QtSizePolicyRigid = QSizePolicy::Fixed;
const int QtSizePolicyControlled = QSizePolicy::Ignored;
const int QtSizePolicyIncompressible = QSizePolicy::Minimum;
const int QtSizePolicyIncompressibleExpanding = QSizePolicy::MinimumExpanding;
const int QtSizePolicyFree = QSizePolicy::Preferred;
const int QtSizePolicyFreeExpanding = QSizePolicy::Expanding;
const int QtSizePolicyBounded = QSizePolicy::Maximum;

const int QtToolButtonIconOnly = Qt::ToolButtonIconOnly;
const int QtToolButtonTextOnly = Qt::ToolButtonTextOnly;
const int QtToolButtonTextBesideIcon = Qt::ToolButtonTextBesideIcon;
const int QtToolButtonTextUnderIcon = Qt::ToolButtonTextUnderIcon;

const int QtInputText = QInputDialog::TextInput;
const int QtInputInt = QInputDialog::IntInput;
const int QtInputDouble = QInputDialog::DoubleInput;

const int QtMsgBoxInfo = QMessageBox::Information;
const int QtMsgBoxWarning = QMessageBox::Warning;
const int QtMsgBoxCritical = QMessageBox::Critical;
const int QtMsgBoxQuestion = QMessageBox::Question;

const int QtFileDialogModeSave = QFileDialog::AnyFile;
const int QtFileDialogModeOpenSingle = QFileDialog::ExistingFile;
const int QtFileDialogModeOpenMultiple = QFileDialog::ExistingFiles;

const int QtMsgBoxOK = QMessageBox::Ok;
const int QtMsgBoxCancel = QMessageBox::Cancel;
const int QtMsgBoxYes = QMessageBox::Yes;
const int QtMsgBoxNo = QMessageBox::No;
const int QtMsgBoxAbort = QMessageBox::Abort;
const int QtMsgBoxRetry = QMessageBox::Retry;
const int QtMsgBoxIgnore = QMessageBox::Ignore;
const int QtMsgBoxSave = QMessageBox::Save;
const int QtMsgBoxDiscard = QMessageBox::Discard;

const int QtBtnBoxOK = QDialogButtonBox::Ok;
const int QtBtnBoxCancel = QDialogButtonBox::Cancel;

const int QtItemNoSelection = QAbstractItemView::NoSelection;
const int QtItemSingleSelection = QAbstractItemView::SingleSelection;
const int QtItemMultiSelection = QAbstractItemView::MultiSelection;
const int QtItemExtendedSelection = QAbstractItemView::ExtendedSelection;

const int QtTextPlain = Qt::PlainText;
const int QtTextHtml = Qt::RichText;
const int QtTextMarkdown = Qt::MarkdownText;

const int QtScrollBothDirection = SmartScrollArea::BothDirection;
const int QtScrollVerticalOnly = SmartScrollArea::VerticalOnly;
const int QtScrollHorizontalOnly = SmartScrollArea::HorizontalOnly;

const int QtLwiPrepend = ListWidgetInterface::Prepend;
const int QtLwiAppend = ListWidgetInterface::Append;
const int QtLwiInsertAbove = ListWidgetInterface::InsertAbove;
const int QtLwiInsertBelow = ListWidgetInterface::InsertBelow;

const int QtLwiUp = ListWidgetInterface::Up;
const int QtLwiDown = ListWidgetInterface::Down;

static QApplication*
    app = nullptr;
static CallbackExecutor*
    executor = nullptr;
static bool
    initialized = false;

void QtInit() {
    static int fake_argc = 1;
    static char fake_arg[] = {'Q','t','A','p','p','\0'};
    static char* fake_argv[] = { fake_arg };
    if (!(initialized)) {
        QCoreApplication::setAttribute(Qt::AA_ShareOpenGLContexts);
        app = new QApplication(fake_argc, fake_argv);
        app->setQuitOnLastWindowClosed(false);
        QFont f = app->font();
        f.setPixelSize(Get1remPixels());
        app->setFont(f);
        qRegisterMetaType<callback_t>();
        executor = new CallbackExecutor();
        initialized = true;
    }
}
int QtMain() {
    return app->exec();
}
void QtSchedule(void (*cb)(uint64_t), uint64_t payload) {
    emit executor->QueueCallback(cb, payload);
}
void QtExit(int code) {
    app->exit(code);
}
void QtQuit() {
    app->quit();
}

QtString QtNewUUID() {
    return WrapString(QUuid::createUuid().toString());
}
int QtFontSize() {
    return Get1remPixels();
}

void* QtObjectFindChild(void* object_ptr, const char* name) {
    QObject* obj = (QObject*) object_ptr;
    QObject* child = obj->findChild<QObject*>(QString(name));
    return (void*) child;
}
void* QtWidgetFindChildWidget(void* widget_ptr, const char* name) {
    QWidget* widget = (QWidget*) widget_ptr;
    QWidget* child = widget->findChild<QWidget*>(QString(name));
    return (void*) child;
}
void* QtWidgetFindChildAction(void* widget_ptr, const char* name) {
    QWidget* widget = (QWidget*) widget_ptr;
    QAction* child = widget->findChild<QAction*>(QString(name));
    return (void*) child;
}

void QtWidgetShow(void* widget_ptr) {
    QWidget* widget = (QWidget*) widget_ptr;
    widget->show();
}
void QtWidgetHide(void* widget_ptr) {
    QWidget* widget = (QWidget*) widget_ptr;
    widget->hide();
}
void QtWidgetRaise(void* widget_ptr) {
    QWidget* widget = (QWidget*) widget_ptr;
    widget->raise();
}
void QtWidgetActivateWindow(void* widget_ptr) {
    QWidget* widget = (QWidget*) widget_ptr;
    widget->activateWindow();
    #ifdef _WIN32
    HWND hwnd = (HWND) widget->winId();
    SetWindowPos(hwnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOACTIVATE);
    SetWindowPos(hwnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOMOVE | SWP_NOSIZE | SWP_NOACTIVATE);
    #endif
}
void QtWidgetMoveToScreenCenter(void* widget_ptr) {
    QWidget* widget = (QWidget*) widget_ptr;
    MoveToScreenCenter(widget);
}
void QtWidgetClearTextLater(void* widget_ptr) {
    if (QLineEdit* edit = qobject_cast<QLineEdit*>((QWidget*) widget_ptr)) {
        QTimer::singleShot(0, edit, &QLineEdit::clear);
    }
    if (QPlainTextEdit* edit = qobject_cast<QPlainTextEdit*>((QWidget*) widget_ptr)) {
        QTimer::singleShot(0, edit, &QPlainTextEdit::clear);
    }
}

QtString QtObjectGetClassName(void* obj_ptr) {
    QObject* obj = (QObject*) obj_ptr;
    return WrapString(obj->metaObject()->className());
}
QtBool QtObjectSetPropBool(void* obj_ptr, const char* prop, QtBool val) {
    QObject* obj = (QObject*) obj_ptr;
    return obj->setProperty(prop, (val != 0));
}
QtBool QtObjectGetPropBool(void* obj_ptr, const char* prop) {
    QObject* obj = (QObject*) obj_ptr;
    QVariant val = obj->property(prop);
    return val.toBool();
}
QtBool QtObjectSetPropString(void* obj_ptr, const char* prop, QtString val) {
    QObject* obj = (QObject*) obj_ptr;
    return obj->setProperty(prop, UnwrapString(val));
}
QtString QtObjectGetPropString(void* obj_ptr, const char* prop) {
    QObject* obj = (QObject*) obj_ptr;
    QVariant val = obj->property(prop);
    return WrapString(val.toString());
}
QtBool QtObjectSetPropInt(void* obj_ptr, const char* prop, int val) {
    QObject* obj = (QObject*) obj_ptr;
    return obj->setProperty(prop, val);
}
int QtObjectGetPropInt(void* obj_ptr, const char* prop) {
    QObject* obj = (QObject*) obj_ptr;
    QVariant val = obj->property(prop);
    return val.toInt();
}
QtBool QtObjectSetPropDouble(void* obj_ptr, const char* prop, double val) {
    QObject* obj = (QObject*) obj_ptr;
    return obj->setProperty(prop, val);
}
double QtObjectGetPropDouble(void* obj_ptr, const char* prop) {
    QObject* obj = (QObject*) obj_ptr;
    QVariant val = obj->property(prop);
    return val.toDouble();
}
QtBool QtObjectSetPropPixmap(void* obj_ptr, const char* prop, QtPixmap val) {
    QObject* obj = (QObject*) obj_ptr;
    return obj->setProperty(prop, QVariant(*(QPixmap*)(val.ptr)));
}

void QtDeleteObjectLater(void* obj_ptr) {
    QObject* obj = (QObject*) obj_ptr;
    obj->deleteLater();
}

void* QtConnect (
    void* obj_ptr,
    const char* signal,
    void (*cb)(uint64_t),
    uint64_t payload
) {
    QObject* target_obj = (QObject*) obj_ptr;
    CallbackObject* cb_obj = new CallbackObject(nullptr, cb, payload);
    if (QtDynamicConnect(target_obj, signal, cb_obj, "Slot()")) {
        return (void*) cb_obj;
    } else {
        delete cb_obj;
        return nullptr;
    }
}
void QtBlockSignals(void* obj_ptr, QtBool block) {
    QObject* obj = (QObject*) obj_ptr;
    obj->blockSignals(bool(block));
}

QtEventListener QtListen (
        void*  obj_ptr,
        int    kind,
        QtBool prevent,
        void     (*cb)(uint64_t),
        uint64_t payload
) {
    QObject* obj = (QObject*) obj_ptr;
    QEvent::Type q_kind = (QEvent::Type) kind;
    EventListener* l = new EventListener(q_kind, prevent, cb, payload);
    obj->installEventFilter(l);
    return { (void*) l };
}
QtEvent QtGetCurrentEvent(QtEventListener listener) {
    EventListener* l = (EventListener*) listener.ptr;
    QtEvent wrapped;
    wrapped.ptr = l->current_event;
    return wrapped;
}
void QtUnlisten(void* obj_ptr, QtEventListener listener) {
    QObject* obj = (QObject*) obj_ptr;
    EventListener* l = (EventListener*) listener.ptr;
    obj->removeEventFilter(l);
    l->deleteLater();
}

size_t QtResizeEventGetWidth(QtEvent ev_) {
    QEvent* ev = (QEvent*) ev_.ptr;
    QResizeEvent* resize = dynamic_cast<QResizeEvent*>(ev);
    return resize->size().width();
}
size_t QtResizeEventGetHeight(QtEvent ev_) {
    QEvent* ev = (QEvent*) ev_.ptr;
    QResizeEvent* resize = dynamic_cast<QResizeEvent*>(ev);
    return resize->size().height();
}

QtString QtDynamicPropertyChangeEventGetPropertyName(QtEvent ev) {
    QDynamicPropertyChangeEvent* change = (QDynamicPropertyChangeEvent*) ev.ptr;
    return WrapString(QString::fromUtf8(change->propertyName()));
}

QtVariant QtCreateVariantInvalid() {
    QVariant* ptr = new QVariant();
    return { (void*) ptr };
}
QtVariant QtCreateVariantInt(int value) {
    QVariant* ptr = new QVariant(value);
    return { (void*) ptr };
}
QtVariant QtCreateVariantDouble(double value) {
    QVariant* ptr = new QVariant(value);
    return { (void*) ptr };
}
QtVariant QtCreateVariantString(QtString value) {
    QVariant* ptr = new QVariant(UnwrapString(value));
    return { (void*) ptr };
}
void QtDeleteVariant(QtVariant v) {
    delete (QVariant*)(v.ptr);
}

QtString QtNewStringUTF8(const uint8_t* buf, size_t len) {
    QString* ptr = new QString();
    *ptr = QString::fromUtf8((const char*)(buf), len);
    return { (void*) ptr };
}
QtString QtNewStringUTF16(const uint16_t* buf, size_t len) {
    QString* ptr = new QString();
    *ptr = QString::fromUtf16((const ushort*)(buf), len);
    return { (void*) ptr };
}
QtString QtNewStringUTF32(const uint32_t* buf, size_t len) {
    QString* ptr = new QString();
    *ptr = QString::fromUcs4(buf, len);
    return { (void*) ptr };
}
void QtDeleteString(QtString str) {
    delete (QString*)(str.ptr);
}

size_t QtStringUTF16Length(QtString str) {
    return UnwrapString(str).length();
}
void QtStringWriteToUTF16Buffer(QtString str, uint16_t* buf) {
    QString s = UnwrapString(str);
    for (QChar c: s) {
        *buf = c.unicode();
        buf += 1;
    }
}
size_t QtStringWriteToUTF32Buffer(QtString str, uint32_t* buf) {
    const QVector<uint> vec = UnwrapString(str).toUcs4();
    size_t len = 0;
    for (uint rune: vec) {
        *buf = rune;
        buf += 1;
        len += 1;
    }
    return len;
}

size_t QtStringListGetSize(QtStringList list) {
    QStringList* ptr = (QStringList*) (list.ptr);
    return ptr->size();
}
QtString QtStringListGetItem(QtStringList list, size_t index) {
    QStringList* ptr = (QStringList*) (list.ptr);
    return WrapString(ptr->at(index));
}
void QtDeleteStringList(QtStringList list) {
    delete (QStringList*)(list.ptr);
}

uint8_t* QtByteArrayGetBuffer(QtByteArray data) {
    QByteArray* ptr = (QByteArray*) (data.ptr);
    return (uint8_t*) ptr->data();
}
size_t QtByteArrayGetSize(QtByteArray data) {
    QByteArray* ptr = (QByteArray*) (data.ptr);
    return ptr->size();
}
void QtDeleteByteArray(QtByteArray data) {
    delete (QByteArray*)(data.ptr);
}

QtVariantList QtNewVariantList() {
    return { (void*) new QVariantList() };
}
void QtVariantListAppendNumber(QtVariantList l, double n) {
    QVariantList* ptr = (QVariantList*) l.ptr;
    ptr->append(n);
}
void QtVariantListAppendString(QtVariantList l, QtString str) {
    QVariantList* ptr = (QVariantList*) l.ptr;
    ptr->append(UnwrapString(str));
}
void QtDeleteVariantList(QtVariantList l) {
    delete (QVariantList*)(l.ptr);
}

QtString QtVariantMapGetString(QtVariantMap m, QtString key) {
    QVariantMap* ptr = (QVariantMap*) m.ptr;
    QString key_ = UnwrapString(key);
    QVariant val_ = (*ptr)[key_];
    QtString val = WrapString(val_.toString());
    return val;
}
double QtVariantMapGetFloat(QtVariantMap m, QtString key) {
    QVariantMap* ptr = (QVariantMap*) m.ptr;
    QString key_ = UnwrapString(key);
    QVariant val_ = (*ptr)[key_];
    double val = val_.toDouble();
    return val;
}
QtBool QtVariantMapGetBool(QtVariantMap m, QtString key) {
    QVariantMap* ptr = (QVariantMap*) m.ptr;
    QString key_ = UnwrapString(key);
    QVariant val_ = (*ptr)[key_];
    int val = val_.toBool();
    return val;
}
void QtDeleteVariantMap(QtVariantMap m) {
    delete (QVariantMap*)(m.ptr);
}

QtIcon QtCreateNullIcon() {
    QIcon* ptr = new QIcon();
    return { (void*) ptr };
}
QtIcon QtCreateIconFromStock(QtString name_) {
    QString name = UnwrapString(name_);
    QString path;
    if (name == "qt-logo") {
        path = ":/qtbinding/qt.png";
    } else {
        path = QString(":/Tango/%1.svg").arg(name);
    }
    QIcon* ptr = new QIcon(path);
    return { (void*) ptr };
}
QtIcon QtCreateIconFromFile(QtString path_) {
    QString path = UnwrapString(path_);
    QIcon* ptr = new QIcon(path);
    return { (void*) ptr };
}
void QtDeleteIcon(QtIcon icon) {
    delete (QIcon*)(icon.ptr);
}

QtPixmap QtNewPixmap(const uint8_t* buf, size_t len, const char* format) {
    QPixmap *ptr = new QPixmap;
    ptr->loadFromData(buf, len, format);
    return { (void*) ptr };
}
QtPixmap QtNewPixmapPNG(const uint8_t* buf, size_t len) {
    return QtNewPixmap(buf, len, "PNG");
}
QtPixmap QtNewPixmapJPEG(const uint8_t* buf, size_t len) {
    return QtNewPixmap(buf, len, "JPG");
}
void QtDeletePixmap(QtPixmap pm) {
    delete (QPixmap*)(pm.ptr);
}

void* QtCreateAction(QtIcon icon, QtString text, QtString* shortcuts_ptr, size_t shortcuts_len, QtBool repeat) {
    QAction* a = new QAction(UnwrapIcon(icon), UnwrapString(text), nullptr);
    QList<QKeySequence> shortcuts;
    for (QString key_seq_str: GetStringList(shortcuts_ptr, shortcuts_len)) {
        shortcuts.append(QKeySequence(key_seq_str));
    }
    a->setShortcuts(shortcuts);
    a->setAutoRepeat(repeat);
    if (shortcuts.length() > 0) {
        a->setToolTip(QString
            ("<p style='white-space:pre'>%1&nbsp;&nbsp;<span style='color:gray;font-size:small;'>%2</span></p>")
                .arg(a->toolTip(), a->shortcut().toString(QKeySequence::NativeText))
        );
    }
    return (void*) a;
}
QtBool QtActionInGroup(void* action_ptr) {
    QAction* action = (QAction*) action_ptr;
    return (action->actionGroup() != nullptr);
}

void* QtCreateActionGroup() {
    QActionGroup* g = new QActionGroup(nullptr);
    return (void*) g;
}
void QtActionGroupAddAction(void* group_ptr, void* action_ptr, int index) {
    QActionGroup* group = (QActionGroup*) group_ptr;
    QAction* action = (QAction*) action_ptr;
    action->setProperty("qtbindingActionIndex", index);
    group->addAction(action);
}
int QtActionGroupGetCheckedActionIndex(void* group_ptr) {
    QActionGroup* group = (QActionGroup*) group_ptr;
    QAction* action = group->checkedAction();
    if (action != nullptr) {
        return action->property("qtbindingActionIndex").toInt();
    } else {
        return -1;
    }
}

void* QtCreateMenu(QtIcon icon, QtString text) {
    QMenu* m = new QMenu(nullptr);
    m->setIcon(UnwrapIcon(icon));
    m->setTitle(UnwrapString(text));
    return (void*) m;
}
void QtMenuAddMenu(void* self_ptr, void* menu_ptr) {
    QMenu* self = (QMenu*) self_ptr;
    QMenu* menu = (QMenu*) menu_ptr;
    self->addMenu(menu);
    new MenuHandle(menu, self);
}
void QtMenuAddAction(void* self_ptr, void* action_ptr) {
    QMenu* self = (QMenu*) self_ptr;
    QAction* action = (QAction*) action_ptr;
    self->addAction(action);
}
void QtMenuAddSeparator(void* self_ptr) {
    QMenu* self = (QMenu*) self_ptr;
    self->addSeparator();
}

void* QtBindContextMenu(void* widget_ptr, void* menu_ptr) {
    QWidget* widget = (QWidget*) widget_ptr;
    QMenu* menu = (QMenu*) menu_ptr;
    ContextMenuBinding* b = ContextMenuBinding::TryCreate(widget, menu);
    if (b == nullptr) {
        return nullptr;
    }
    return (void*) b;
}

void* QtCreateMenuBar() {
    QMenuBar* b = new QMenuBar(nullptr);
    return (void*) b;
}
void QtMenuBarAddMenu(void* self_ptr, void* menu_ptr) {
    QMenuBar* self = (QMenuBar*) self_ptr;
    QMenu* menu = (QMenu*) menu_ptr;
    self->addMenu(menu);
    new MenuHandle(menu, self);
}

void* QtCreateToolBar(int tool_btn_style) {
    QToolBar* b = new QToolBar(nullptr);
    b->setFloatable(false);
    b->setMovable(false);
    b->setToolButtonStyle((Qt::ToolButtonStyle) tool_btn_style);
    int rem = Get1remPixels();
    if (b->toolButtonStyle() == Qt::ToolButtonTextUnderIcon) {
        b->setIconSize(QSize(3*rem,3*rem));
    } else if (b->toolButtonStyle() == Qt::ToolButtonTextBesideIcon) {
        b->setIconSize(QSize(1*rem,1*rem));
    } else {
        b->setIconSize(QSize(2*rem,2*rem));
    }
    return (void*) b;
}
void QtToolBarAddMenu(void* self_ptr, void* menu_ptr) {
    QToolBar* self = (QToolBar*) self_ptr;
    QMenu* menu = (QMenu*) menu_ptr;
    // self->addAction(menu->menuAction()); new MenuHandle(menu, self);
    self->addWidget(menu);
}
void QtToolBarAddAction(void* self_ptr, void* action_ptr) {
    QToolBar* self = (QToolBar*) self_ptr;
    QAction* action = (QAction*) action_ptr;
    QString tooltip;
    self->addAction(action);
    if (self->toolButtonStyle() != Qt::ToolButtonIconOnly
    && action->shortcuts().length() == 0) {
        for(QToolButton* btn: self->findChildren<QToolButton*>()) {
            if (btn->defaultAction() == action) {
                btn->setToolTip("");
                QObject::connect(action, &QAction::changed, btn, [btn] () -> void {
                    btn->setToolTip("");
                });
            }
        }
    }
}
void QtToolBarAddSeparator(void* self_ptr) {
    QToolBar* self = (QToolBar*) self_ptr;
    self->addSeparator();
}
void QtToolBarAddWidget(void* self_ptr, void* widget_ptr) {
    QToolBar* self = (QToolBar*) self_ptr;
    QWidget* widget = (QWidget*) widget_ptr;
    self->addWidget(WeakWidget::Wrap(widget));
}
void QtToolBarAddSpacer(void* self_ptr, int width, int height, QtBool expand) {
    QToolBar* self = (QToolBar*) self_ptr;
    width = GetScaledLength(width);
    height = GetScaledLength(height);
    QSizePolicy::Policy policy = (expand)? QSizePolicy::MinimumExpanding: QSizePolicy::Minimum;
    SizeHintWidget* spacer = new SizeHintWidget(QSize(width, height), nullptr);
    if (self->orientation() == Qt::Horizontal) {
        spacer->setSizePolicy(policy, QSizePolicy::Minimum);
    } else {
        spacer->setSizePolicy(QSizePolicy::Minimum, policy);
    }
    self->addWidget(spacer);
}

void* QtCreateDialogButtonBox() {
    QDialogButtonBox* box = new QDialogButtonBox();
    return (void*) box;
}
void* QtDialogButtonBoxAddButton(void* box_ptr, int kind) {
    QDialogButtonBox* box = (QDialogButtonBox*) box_ptr;
    QPushButton* btn = box->addButton((QDialogButtonBox::StandardButton) kind);
    return (void*) btn;
}

void QtDialogAccept(void *dialog_ptr) {
    QDialog* dialog = (QDialog*) dialog_ptr;
    dialog->accept();
}
void QtDialogReject(void *dialog_ptr) {
    QDialog* dialog = (QDialog*) dialog_ptr;
    dialog->reject();
}
int QtDialogGetResult(void* dialog_ptr) {
    QDialog* dialog = (QDialog*) dialog_ptr;
    return dialog->result();
}
QtBool QtDialogGetResultBoolean(void* dialog_ptr) {
    QDialog* dialog = (QDialog*) dialog_ptr;
    return (dialog->result() == QDialog::Accepted);
}

void QtConsumeDialog(void* dialog_ptr, void (*cb)(uint64_t), uint64_t payload) {
    QDialog* dialog = (QDialog*) dialog_ptr;
    if (QApplication::activeWindow() == nullptr) {
        dialog->setWindowFlag(Qt::WindowStaysOnTopHint, true);
    }
    QObject::connect(dialog, &QDialog::finished, [dialog,cb,payload] () -> void {
        cb(payload);
        dialog->deleteLater();
    });
    if (QMessageBox* msgbox = qobject_cast<QMessageBox*>(dialog)) {
        bool missing_reject_btn = true;
        for (QAbstractButton* btn: msgbox->buttons()) {
            if (msgbox->buttonRole(btn) == QMessageBox::RejectRole) {
                missing_reject_btn = false;
            }
        }
        if (missing_reject_btn) {
            QPushButton* cancel = msgbox->addButton(QMessageBox::Cancel);
            cancel->setVisible(false);
        }
    }
    dialog->show();
    if (QMessageBox* msgbox = qobject_cast<QMessageBox*>(dialog)) {
        msgbox->defaultButton()->setFocus();
    }
}

void* QtCreateInputDialog(int mode, QtVariant value, QtString title, QtString prompt) {
    QInputDialog* dialog = new QInputDialog(nullptr);
    dialog->setWindowModality(Qt::ApplicationModal);
    dialog->setInputMode((QInputDialog::InputMode) mode);
    dialog->setIntMinimum(-2147483648);
    dialog->setIntMaximum(2147483647);
    dialog->setDoubleMinimum(-9007199254740991);
    dialog->setDoubleMaximum(9007199254740991);
    dialog->setDoubleDecimals(3);
    QVariant v = UnwrapVariant(value);
    switch (mode) {
    case QtInputInt:    dialog->setIntValue(v.toInt());       break;
    case QtInputDouble: dialog->setDoubleValue(v.toDouble()); break;
    case QtInputText:   dialog->setTextValue(v.toString());   break;
    }
    dialog->setWindowTitle(UnwrapString(title));
    dialog->setLabelText(UnwrapString(prompt));
    return (void*) dialog;
}
void QtInputDialogUseMultilineText(void* dialog_ptr) {
    QInputDialog* dialog = (QInputDialog*) dialog_ptr;
    dialog->setOption(QInputDialog::UsePlainTextEditForTextInput, true);
}
void QtInputDialogUseChoiceItems(void* dialog_ptr, QtString* items_ptr, size_t items_len) {
    QInputDialog* dialog = (QInputDialog*) dialog_ptr;
    QStringList items;
    for (size_t i = 0; i < items_len; i += 1) {
        items.append(UnwrapString(items_ptr[i]));
    }
    dialog->setComboBoxItems(items);
}
QtString QtInputDialogGetTextValue(void* dialog_ptr) {
    QInputDialog* dialog = (QInputDialog*) dialog_ptr;
    return WrapString(dialog->textValue());
}
int QtInputDialogGetIntValue(void* dialog_ptr) {
    QInputDialog* dialog = (QInputDialog*) dialog_ptr;
    return dialog->intValue();
}
double QtInputDialogGetDoubleValue(void* dialog_ptr) {
    QInputDialog* dialog = (QInputDialog*) dialog_ptr;
    return dialog->doubleValue();
}

void* QtCreateMessageBox(int icon, int buttons, QtString title, QtString content) {
    QMessageBox* msgbox = new QMessageBox(
        (QMessageBox::Icon) icon,
        UnwrapString(title),
        UnwrapString(content),
        (QMessageBox::StandardButtons) buttons,
        nullptr
    );
    msgbox->setWindowModality(Qt::ApplicationModal);
    return (void*) msgbox;
}
void QtMessageBoxSetDefaultButton(void* msgbox_ptr, int btn) {
    QMessageBox* msgbox = (QMessageBox*) msgbox_ptr;
    msgbox->setDefaultButton((QMessageBox::StandardButton) btn);
}
int QtMessageBoxGetResultButton(void* msgbox_ptr) {
    QMessageBox* msgbox = (QMessageBox*) msgbox_ptr;
    return msgbox->result();
}

void* QtCreateFileDialog(int mode_, QtString filters_) {
    QFileDialog::FileMode mode = (QFileDialog::FileMode) mode_;
    QString filters = UnwrapString(filters_);
    QFileDialog* d = new QFileDialog(nullptr);
    d->setWindowModality(Qt::ApplicationModal);
    if (mode_ == QtFileDialogModeSave) {
        d->setAcceptMode(QFileDialog::AcceptSave);
    } else {
        d->setAcceptMode(QFileDialog::AcceptOpen);
    }
    d->setFileMode(mode);
    d->setNameFilter(filters);
    d->setDirectory(QDir::homePath());
    return d;
}
int QtFileDialogGetResultFileCount(void* d_ptr) {
    QFileDialog* d = (QFileDialog*) d_ptr;
    return d->selectedFiles().length();
}
QtString QtFileDialogGetResultFileItem(void* d_ptr, int index) {
    QFileDialog* d = (QFileDialog*) d_ptr;
    return WrapString(d->selectedFiles().at(index));
}

void* QtCreateLayoutRow(int spacing) {
    QBoxLayout* row = new QHBoxLayout();
    row->setSpacing(GetScaledLength(spacing));
    QLayout* layout = (QLayout*) row;
    return (void*) layout;
}
void* QtCreateLayoutColumn(int spacing) {
    QBoxLayout* column = new QVBoxLayout();
    column->setSpacing(GetScaledLength(spacing));
    QLayout* layout = (QLayout*) column;
    return (void*) layout;
}
void* QtCreateLayoutGrid(int row_spacing, int column_spacing) {
    QGridLayout* grid = new QGridLayout();
    grid->setHorizontalSpacing(GetScaledLength(row_spacing));
    grid->setVerticalSpacing(GetScaledLength(column_spacing));
    QLayout* layout = (QLayout*) grid;
    return (void*) layout;
}
QtGridSpan QtMakeGridSpan(int row, int column, int rowSpan, int columnSpan) {
    QtGridSpan span = { row, column, rowSpan, columnSpan };
    return span;
}
void QtLayoutAddLayout(void* self_ptr, void* layout_ptr, QtGridSpan span, int align) {
    QLayout* self = (QLayout*) self_ptr;
    QLayout* layout = (QLayout*) layout_ptr;
    if (QGridLayout* self_grid = qobject_cast<QGridLayout*>(self)) {
        self_grid->addLayout(layout, span.row, span.column, span.rowSpan, span.columnSpan, (Qt::Alignment) align);
    } else {
        self->addItem(layout);
    }
}
void QtLayoutAddWidget(void* self_ptr, void* widget_ptr, QtGridSpan span, int align) {
    QLayout* self = (QLayout*) self_ptr;
    QWidget* widget = (QWidget*) widget_ptr;
    if (widget->parentWidget() != nullptr) {
        return;
    }
    if (QGridLayout* self_grid = qobject_cast<QGridLayout*>(self)) {
        self_grid->addWidget(widget, span.row, span.column, span.rowSpan, span.columnSpan, (Qt::Alignment) align);
    } else {
        self->addWidget(widget);
    }
    QObject::connect(widget, &QObject::destroyed, self, [widget,self] () -> void {
        self->removeWidget(widget);
    });
}
void QtLayoutAddSpacer(void* self_ptr, int width, int height, QtBool expand, QtGridSpan span, int align) {
    QLayout* self = (QLayout*) self_ptr;
    width = GetScaledLength(width);
    height = GetScaledLength(height);
    QSizePolicy::Policy policy = (expand)? QSizePolicy::MinimumExpanding: QSizePolicy::Minimum;
    if (QGridLayout* self_grid = qobject_cast<QGridLayout*>(self)) {
        QSpacerItem* spacer_item = new QSpacerItem(width, height, policy, policy);
        self_grid->addItem(spacer_item, span.row, span.column, span.rowSpan, span.columnSpan, (Qt::Alignment) align);
    } else if (QBoxLayout* self_box = qobject_cast<QBoxLayout*>(self)) {
        QBoxLayout::Direction d = self_box->direction();
        if (d == QBoxLayout::LeftToRight || d == QBoxLayout::RightToLeft) {
            QSpacerItem* spacer_item = new QSpacerItem(width, height, policy, QSizePolicy::Minimum);
            self->addItem(spacer_item);
        } else {
            QSpacerItem* spacer_item = new QSpacerItem(width, height, QSizePolicy::Minimum, policy);
            self->addItem(spacer_item);
        }
    } else {
        QSpacerItem* spacer_item = new QSpacerItem(width, height, policy, policy);
        self->addItem(spacer_item);
    }
}
void QtLayoutAddLabel(void* self_ptr, QtString text, QtGridSpan span, int align) {
    QLayout* self = (QLayout*) self_ptr;
    QLabel* label = new QLabel(UnwrapString(text), nullptr);
    if (QGridLayout* self_grid = qobject_cast<QGridLayout*>(self)) {
        self_grid->addWidget(label, span.row, span.column, span.rowSpan, span.columnSpan, (Qt::Alignment) align);
    } else {
        self->addWidget(label);
    }
    QObject::connect(self, &QObject::destroyed, label, &QObject::deleteLater);
}

QLayout* GetLayoutWithMargins(void* layout_ptr, int margin_x, int margin_y) {
    QLayout* layout = (QLayout*) layout_ptr;
    margin_x = GetScaledLength(margin_x);
    margin_y = GetScaledLength(margin_y);
    layout->setContentsMargins(margin_x, margin_y, margin_x, margin_y);
    return layout;
}
bool SetSize(QWidget* w, int width, int height) {
    if (width > 0 && height > 0) {
        width = GetScaledLength(width);
        height = GetScaledLength(height);
        w->resize(width, height);
        return true;
    } else {
        return false;
    }
}
void* QtCreateWidget(void* layout_ptr, int margin_x, int margin_y, int policy_x, int policy_y) {
    QLayout* layout = GetLayoutWithMargins(layout_ptr, margin_x, margin_y);
    WeakWidget* widget = new WeakWidget(layout);
    widget->setSizePolicy((QSizePolicy::Policy) policy_x, (QSizePolicy::Policy) policy_y);
    return (void*) widget;
}
void* QtCreateMainWindow(void* menu_bar_ptr, void* tool_bar_ptr, void* layout_ptr, int margin_x, int margin_y, int width, int height, QtIcon icon) {
    QLayout* layout = GetLayoutWithMargins(layout_ptr, margin_x, margin_y);
    QMainWindow* window = new QMainWindow(nullptr);
    SetSize(window, width, height);
    window->setWindowIcon(UnwrapIcon(icon));
    window->setCentralWidget(new WeakWidget(layout));
    if (QMenuBar* menu_bar = (QMenuBar*) menu_bar_ptr) {
        window->setMenuBar(menu_bar);
    }
    if (QToolBar* tool_bar = (QToolBar*) tool_bar_ptr) {
        window->addToolBar(tool_bar);
    }
    return (void*) window;
}
void* QtCreateDialog(void* layout_ptr, int margin_x, int margin_y, int width, int height, QtIcon icon) {
    QLayout* layout = GetLayoutWithMargins(layout_ptr, margin_x, margin_y);
    QDialog* dialog = (QDialog*) new CustomDialog();
    dialog->setWindowModality(Qt::ApplicationModal);
    bool fixed = !(SetSize(dialog, width, height));
    dialog->setWindowIcon(UnwrapIcon(icon));
    QVBoxLayout* wrapper = new QVBoxLayout();
    wrapper->addWidget(new WeakWidget(layout));
    wrapper->setContentsMargins(0, 0, 0, 0);
    if (fixed) {
        wrapper->setSizeConstraint(QLayout::SetFixedSize);
    }
    dialog->setLayout(wrapper);
    return (void*) dialog;
}
void* QtCreateScrollArea(int direction, void* layout_ptr, int margin_x, int margin_y) {
    QLayout* layout = GetLayoutWithMargins(layout_ptr, margin_x, margin_y);
    auto d = (SmartScrollArea::Direction) direction;
    SmartScrollArea* area = new SmartScrollArea(d, nullptr);
    area->setWidget(new WeakWidget(layout));
    return (void*) area;
}
void* QtCreateGroupBox(QtString title, void* layout_ptr, int margin_x, int margin_y) {
    QLayout* layout = GetLayoutWithMargins(layout_ptr, margin_x, margin_y);
    QGroupBox* group = new QGroupBox(UnwrapString(title), nullptr);
    QVBoxLayout* wrapper = new QVBoxLayout();
    wrapper->addWidget(new WeakWidget(layout));
    group->setLayout(wrapper);
    return (void*) group;
}
void* QtCreateSplitter(void** widgets_ptr, size_t widgets_len) {
    QWidgetList widgets = GetWidgetList(widgets_ptr, widgets_len);
    QSplitter* splitter = new QSplitter();
    for (QWidget* widget: widgets) {
        splitter->addWidget(WeakWidget::Wrap(widget));
    }
    splitter->setChildrenCollapsible(false);
    return (void*) splitter;
}

void* QtCreateDynamicWidget() {
    return (void*) new DynamicWidget(nullptr);
}
void QtDynamicWidgetSetWidget(void* self_ptr, void* widget_ptr) {
    DynamicWidget* self = (DynamicWidget*) self_ptr;
    QWidget* w = (QWidget*) widget_ptr;
    self->setWidget(w);
}

void* QtCreateDummyFocusableWidget() {
    QWidget* w = new QWidget();
    w->setFocusPolicy(Qt::StrongFocus);
    return (void*) w;
}

void* QtCreateLabel(QtString text, int align, QtBool selectable) {
    QLabel* l = new QLabel(UnwrapString(text), nullptr);
    l->setTextFormat(Qt::PlainText);
    l->setAlignment((Qt::Alignment) align);
    if (selectable) {
        l->setTextInteractionFlags(Qt::TextSelectableByMouse | l->textInteractionFlags());
    }
    return (void*) l;
}
void* QtCreateIconLabel(QtIcon icon_, int size) {
    QLabel* l = new QLabel(nullptr);
    QIcon icon = UnwrapIcon(icon_);
    if (size <= 0) {
        bool ok = false;
        for (QSize size: icon.availableSizes()) {
            l->setPixmap(icon.pixmap(size));
            ok = true;
            break;
        }
        if (!(ok)) {
            l->setPixmap(icon.pixmap(Get1remPixels()));
        }
    } else {
        l->setPixmap(icon.pixmap((size * Get1remPixels())));
    }
    return (void*) l;
}
void* QtCreateElidedLabel(QtString text) {
    ElidedLabel* l = new ElidedLabel(UnwrapString(text), nullptr);
    return (void*) l;
}
void* QtCreateTextView(QtString text, int format) {
    TextView* v = new TextView(UnwrapString(text), (Qt::TextFormat) format, nullptr);
    return (void*) v;
}
void* QtCreateCheckBox(QtString text, QtBool checked) {
    QCheckBox* cb = new QCheckBox(UnwrapString(text));
    cb->setChecked(checked);
    return (void*) cb;
}
void* QtCreateComboBox() {
    QComboBox* cb = new QComboBox();
    return (void*) cb;
}
void QtComboBoxAddItem(void* b_ptr, QtIcon icon, QtString name, QtBool selected) {
    QComboBox* cb = (QComboBox*) b_ptr;
    int index = cb->count();
    cb->addItem(UnwrapIcon(icon), UnwrapString(name));
    if (selected) {
        cb->setCurrentIndex(index);
    }
}
void* QtCreateComboBoxDialog(QtString title, QtString prompt) {
    ComboBoxDialog* d = new ComboBoxDialog(UnwrapString(prompt), nullptr);
    d->setWindowModality(Qt::ApplicationModal);
    d->setWindowTitle(UnwrapString(title));
    return (void*) d;
}
void* QtComboBoxDialogGetComboBox(void* d_ptr) {
    ComboBoxDialog* d = (ComboBoxDialog*) d_ptr;
    QComboBox* b = d->ComboBox();
    return (void*) b;
}
void* QtCreatePushButton(QtIcon icon, QtString text, QtString tooltip) {
    QPushButton* b = new QPushButton(UnwrapIcon(icon), UnwrapString(text));
    b->setToolTip(UnwrapString(tooltip));
    int rem = Get1remPixels();
    b->setIconSize(QSize(1*rem,1*rem));
    return (void*) b;
}
void* QtCreateLineEdit(QtString text) {
    return (void*) (new QLineEdit(UnwrapString(text), nullptr));
}
void* QtCreatePlainTextEdit(QtString text) {
    return (void*) (new QPlainTextEdit(UnwrapString(text), nullptr));
}
void* QtCreateSlider(int min, int max, int value) {
    QSlider* s = new QSlider(Qt::Horizontal, nullptr);
    s->setMinimum(min);
    s->setMaximum(max);
    s->setSingleStep(1);
    s->setPageStep(1);
    s->setValue(value);
    return (void*) s;
}
void* QtCreateProgressBar(QtString format_, int max) {
    QProgressBar* pb = new QProgressBar(nullptr);
    QString format = UnwrapString(format_);
    if (format != "") {
        pb->setTextVisible(true);
        pb->setFormat(format);
    } else {
        pb->setTextVisible(false);
    }
    pb->setMinimum(0);
    pb->setMaximum(max);
    pb->setValue(0);
    return (void*) pb;
}

QtString QtClipboardReadText() {
    QClipboard* c = QGuiApplication::clipboard();
    return WrapString(c->text());
}
void QtClipboardWriteText(QtString text) {
    QClipboard* c = QGuiApplication::clipboard();
    c->setText(UnwrapString(text));
}

void* QtLwiCreateFromDefaultListWidget(size_t columns, int select_, void** headers_ptr, size_t headers_len, int stretch) {
    auto pick = (DefaultListWidget::Select) select_;
    auto headers = GetWidgetList(headers_ptr, headers_len);
    auto w = new DefaultListWidget(columns, pick, headers, stretch, nullptr);
    auto lwi = static_cast<ListWidgetInterface*>(w);
    return (void*) lwi;
}
void* QtLwiCastToWidget(void* lwi_ptr) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto w = lwi->CastToWidget();
    return (void*) w;
}
QtString QtLwiCurrent(void* lwi_ptr, QtBool* exists) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    return WrapString(lwi->Current((bool*) exists));
}
QtStringList QtLwiAll(void* lwi_ptr) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto ptr = new QStringList();
    *ptr = lwi->All();
    return { ptr };
}
QtStringList QtLwiSelection(void* lwi_ptr) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto ptr = new QStringList();
    *ptr = lwi->Selection();
    return { ptr };
}
QtStringList QtLwiContiguousSelection(void* lwi_ptr) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto ptr = new QStringList();
    *ptr = lwi->ContiguousSelection();
    return { ptr };
}
void QtLwiInsert(void* lwi_ptr, int hint_, QtString pivot_, QtString key_, void** widgets_ptr, size_t widgets_len) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto hint = (ListWidgetInterface::InsertHint) hint_;
    auto pivot = UnwrapString(pivot_);
    auto key = UnwrapString(key_);
    auto widgets = GetWidgetList(widgets_ptr, widgets_len);
    lwi->Insert(hint, pivot, key, widgets);
}
void QtLwiUpdate(void* lwi_ptr, QtString key_) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto key = UnwrapString(key_);
    lwi->Update(key);
}
QtBool QtLwiMove(void* lwi_ptr, int hint_, QtString key_) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto hint = (ListWidgetInterface::MoveHint) hint_;
    auto key = UnwrapString(key_);
    return lwi->Move(hint, key);
}
void QtLwiDelete(void* lwi_ptr, QtString key_) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto key = UnwrapString(key_);
    lwi->Delete(key);
}
void QtLwiReorder(void* lwi_ptr, QtString* order_ptr, size_t order_len) {
    auto lwi = (ListWidgetInterface*) lwi_ptr;
    auto order = GetStringList(order_ptr, order_len);
    lwi->Reorder(order);
}


