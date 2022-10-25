#ifndef QTBINDING_H
#define QTBINDING_H

#include <stdlib.h>
#include <stdint.h>


#ifdef _WIN32
	#ifdef QTBINDING_WIN32_DLL
		#define EXPORT __declspec(dllexport)
	#else
		#define EXPORT __declspec(dllimport)
	#endif
#else
	#define EXPORT
#endif

typedef int QtBool;

struct _QtString {
    void* ptr;
};
typedef struct _QtString QtString;

struct _QtStringList {
    void* ptr;
};
typedef struct _QtStringList QtStringList;

struct _QtVariant {
    void* ptr;
};
typedef struct _QtVariant QtVariant;

struct _QtVariantMap {
    void* ptr;
};
typedef struct _QtVariantMap QtVariantMap;

struct _QtVariantList {
    void* ptr;
};
typedef struct _QtVariantList QtVariantList;

struct _QtIcon {
    void* ptr;
};
typedef struct _QtIcon QtIcon;

struct _QtPixmap {
    void* ptr;
};
typedef struct _QtPixmap QtPixmap;

struct _QtEvent {
    void* ptr;
};
typedef struct _QtEvent QtEvent;

struct _QtEventListener {
    void* ptr;
};
typedef struct _QtEventListener QtEventListener;

struct _QtByteArray {
    void* ptr;
};
typedef struct _QtByteArray QtByteArray;

struct _QtGridSpan {
    int row;
    int column;
    int rowSpan;
    int columnSpan;
};
typedef struct _QtGridSpan QtGridSpan;


#ifdef __cplusplus
extern "C" {
#endif
	// Event Categories
	EXPORT extern const int QtEventMove;
	EXPORT extern const int QtEventResize;
    EXPORT extern const int QtEventShow;
	EXPORT extern const int QtEventClose;
    EXPORT extern const int QtEventFocusIn;
    EXPORT extern const int QtEventFocusOut;
    EXPORT extern const int QtEventWindowActivate;
    EXPORT extern const int QtEventWindowDeactivate;
	EXPORT extern const int QtEventDynamicPropertyChange;
    // Alignment
    EXPORT extern const int QtAlignDefault;
    EXPORT extern const int QtAlignLeft;
    EXPORT extern const int QtAlignRight;
    EXPORT extern const int QtAlignHCenter;
    EXPORT extern const int QtAlignTop;
    EXPORT extern const int QtAlignBottom;
    EXPORT extern const int QtAlignVCenter;
    // Size Policy
    EXPORT extern const int QtSizePolicyRigid;
    EXPORT extern const int QtSizePolicyControlled;
    EXPORT extern const int QtSizePolicyIncompressible;
    EXPORT extern const int QtSizePolicyIncompressibleExpanding;
    EXPORT extern const int QtSizePolicyFree;
    EXPORT extern const int QtSizePolicyFreeExpanding;
    EXPORT extern const int QtSizePolicyBounded;
    // Tool Button Style
    EXPORT extern const int QtToolButtonIconOnly;
    EXPORT extern const int QtToolButtonTextOnly;
    EXPORT extern const int QtToolButtonTextBesideIcon; 
    EXPORT extern const int QtToolButtonTextUnderIcon;
    // Input Dialog Modes
    EXPORT extern const int QtInputText;
    EXPORT extern const int QtInputInt;
    EXPORT extern const int QtInputDouble;
    // Message Box Icons
    EXPORT extern const int QtMsgBoxInfo;
    EXPORT extern const int QtMsgBoxWarning;
    EXPORT extern const int QtMsgBoxCritical;
    EXPORT extern const int QtMsgBoxQuestion;
    // File Dialog Modes
    EXPORT extern const int QtFileDialogModeSave;
    EXPORT extern const int QtFileDialogModeOpenSingle; 
    EXPORT extern const int QtFileDialogModeOpenMultiple;
    // Message Box Standard Buttons
    EXPORT extern const int QtMsgBoxOK;
    EXPORT extern const int QtMsgBoxCancel;
    EXPORT extern const int QtMsgBoxYes;
    EXPORT extern const int QtMsgBoxNo;
    EXPORT extern const int QtMsgBoxAbort;
    EXPORT extern const int QtMsgBoxRetry;
    EXPORT extern const int QtMsgBoxIgnore;
    EXPORT extern const int QtMsgBoxSave;
    EXPORT extern const int QtMsgBoxDiscard;
    // Dialog Button Box Standard Buttons
    EXPORT extern const int QtBtnBoxOK;
    EXPORT extern const int QtBtnBoxCancel;
    // Item Selection Mode
    EXPORT extern const int QtItemNoSelection;
    EXPORT extern const int QtItemSingleSelection;
    EXPORT extern const int QtItemMultiSelection;
    EXPORT extern const int QtItemExtendedSelection;
    // Text Format
    EXPORT extern const int QtTextPlain;
    EXPORT extern const int QtTextHtml;
    EXPORT extern const int QtTextMarkdown;
    // Scroll Direction
    EXPORT extern const int QtScrollBothDirection;
    EXPORT extern const int QtScrollVerticalOnly;
    EXPORT extern const int QtScrollHorizontalOnly;
    // ListWidgetInterface InsertHint
    EXPORT extern const int QtLwiPrepend;
    EXPORT extern const int QtLwiAppend;
    EXPORT extern const int QtLwiInsertAbove;
    EXPORT extern const int QtLwiInsertBelow;
    // ListWidgetInterface MoveHint
    EXPORT extern const int QtLwiUp;
    EXPORT extern const int QtLwiDown;
	//
    EXPORT void QtInit();
    EXPORT int QtMain();
    EXPORT void QtSchedule(void (*cb)(uint64_t), uint64_t payload);
    EXPORT void QtExit(int code);
    EXPORT void QtQuit();
    EXPORT QtString QtNewUUID();
    EXPORT int QtFontSize();
    EXPORT void* QtObjectFindChild(void* object_ptr, const char* name);
    EXPORT void* QtWidgetFindChildWidget(void* widget_ptr, const char* name);
    EXPORT void* QtWidgetFindChildAction(void* widget_ptr, const char* name);
    EXPORT void QtWidgetShow(void* widget_ptr);
    EXPORT void QtWidgetHide(void* widget_ptr);
    EXPORT void QtWidgetRaise(void* widget_ptr);
    EXPORT void QtWidgetActivateWindow(void* widget_ptr);
    EXPORT void QtWidgetMoveToScreenCenter(void* widget_ptr);
    EXPORT void QtWidgetClearTextLater(void* widget_ptr);
    EXPORT QtString QtObjectGetClassName(void* obj_ptr);
    EXPORT QtBool QtObjectSetPropBool(void* obj_ptr, const char* prop, QtBool val);
    EXPORT QtBool QtObjectGetPropBool(void* obj_ptr, const char* prop);
    EXPORT QtBool QtObjectSetPropString(void* obj_ptr, const char* prop, QtString val);
    EXPORT QtString QtObjectGetPropString(void* obj_ptr, const char* prop);
    EXPORT QtBool QtObjectSetPropInt(void* obj_ptr, const char* prop, int val);
    EXPORT int QtObjectGetPropInt(void* obj_ptr, const char* prop);
    EXPORT QtBool QtObjectSetPropDouble(void* obj_ptr, const char* prop, double val);
    EXPORT double QtObjectGetPropDouble(void* obj_ptr, const char* prop);
    EXPORT QtBool QtObjectSetPropPixmap(void* obj_ptr, const char* prop, QtPixmap val);
    EXPORT void QtDeleteObjectLater(void* obj_ptr);
    EXPORT void* QtConnect(void* obj_ptr, const char* signal, void (*cb)(uint64_t), uint64_t payload);
    EXPORT void QtBlockSignals(void* obj_ptr, QtBool block);
    EXPORT QtEventListener QtListen(void* obj_ptr, int kind, QtBool prevent, void (*cb)(uint64_t), uint64_t payload);
    EXPORT QtEvent QtGetCurrentEvent(QtEventListener listener);
    EXPORT void QtUnlisten(void* obj_ptr, QtEventListener listener);
    EXPORT size_t QtResizeEventGetWidth(QtEvent ev);
    EXPORT size_t QtResizeEventGetHeight(QtEvent ev);
    EXPORT QtString QtDynamicPropertyChangeEventGetPropertyName(QtEvent ev);
    EXPORT QtVariant QtCreateVariantInvalid();
    EXPORT QtVariant QtCreateVariantInt(int value);
    EXPORT QtVariant QtCreateVariantDouble(double value);
    EXPORT QtVariant QtCreateVariantString(QtString value);
    EXPORT void QtDeleteVariant(QtVariant v);
    EXPORT QtString QtNewStringUTF8(const uint8_t* buf, size_t len);
    EXPORT QtString QtNewStringUTF16(const uint16_t* buf, size_t len);
    EXPORT QtString QtNewStringUTF32(const uint32_t* buf, size_t len);
    EXPORT void QtDeleteString(QtString str);
    EXPORT size_t QtStringListGetSize(QtStringList list);
    EXPORT QtString QtStringListGetItem(QtStringList list, size_t index);
    EXPORT void QtDeleteStringList(QtStringList list);
    EXPORT uint8_t* QtByteArrayGetBuffer(QtByteArray data);
    EXPORT size_t QtByteArrayGetSize(QtByteArray data);
    EXPORT void QtDeleteByteArray(QtByteArray data);
    EXPORT QtVariantList QtNewVariantList();
    EXPORT void QtVariantListAppendNumber(QtVariantList l, double n);
    EXPORT void QtVariantListAppendString(QtVariantList l, QtString str);
    EXPORT void QtDeleteVariantList(QtVariantList l);
    EXPORT QtString QtVariantMapGetString(QtVariantMap m, QtString key);
    EXPORT double QtVariantMapGetFloat(QtVariantMap m, QtString key);
    EXPORT QtBool QtVariantMapGetBool(QtVariantMap m, QtString key);
    EXPORT void QtDeleteVariantMap(QtVariantMap m);
    EXPORT size_t QtStringUTF16Length(QtString str);
    EXPORT void QtStringWriteToUTF16Buffer(QtString str, uint16_t* buf);
    EXPORT size_t QtStringWriteToUTF32Buffer(QtString str, uint32_t *buf);
    EXPORT QtIcon QtCreateNullIcon();
    EXPORT QtIcon QtCreateIconFromStock(QtString name_);
    EXPORT QtIcon QtCreateIconFromFile(QtString path_);
    EXPORT QtIcon QtNewIconFromPixmap(QtPixmap pm);
    EXPORT void QtDeleteIcon(QtIcon icon);
    EXPORT QtPixmap QtNewPixmapPNG(const uint8_t* buf, size_t len);
    EXPORT QtPixmap QtNewPixmapJPEG(const uint8_t* buf, size_t len);
    EXPORT void QtDeletePixmap(QtPixmap pm);
    EXPORT void* QtCreateAction(QtIcon icon, QtString text, QtString* shortcuts_ptr, size_t shortcuts_len, QtBool repeat);
    EXPORT QtBool QtActionInGroup(void* action_ptr);
    EXPORT void* QtCreateActionGroup();
    EXPORT void QtActionGroupAddAction(void* group_ptr, void* action_ptr, int index);
    EXPORT int QtActionGroupGetCheckedActionIndex(void* group_ptr);
    EXPORT void* QtCreateMenu(QtIcon icon, QtString text);
    EXPORT void QtMenuAddMenu(void* self_ptr, void* menu_ptr);
    EXPORT void QtMenuAddAction(void* self_ptr, void* action_ptr);
    EXPORT void QtMenuAddSeparator(void* self_ptr);
    EXPORT void* QtBindContextMenu(void* widget_ptr, void* menu_ptr);
    EXPORT void* QtCreateMenuBar();
    EXPORT void QtMenuBarAddMenu(void* self_ptr, void* menu_ptr);
    EXPORT void* QtCreateToolBar(int tool_btn_style);
    EXPORT void QtToolBarAddMenu(void* self_ptr, void* menu_ptr);
    EXPORT void QtToolBarAddAction(void* self_ptr, void* action_ptr);
    EXPORT void QtToolBarAddSeparator(void* self_ptr);
    EXPORT void QtToolBarAddWidget(void* self_ptr, void* widget_ptr);
    EXPORT void QtToolBarAddSpacer(void* self_ptr, int width, int height, QtBool expand);
    EXPORT void* QtCreateDialogButtonBox();
    EXPORT void* QtDialogButtonBoxAddButton(void* box_ptr, int kind);
    EXPORT void QtDialogAccept(void* dialog_ptr);
    EXPORT void QtDialogReject(void* dialog_ptr);
    EXPORT int QtDialogGetResult(void* dialog_ptr);
    EXPORT QtBool QtDialogGetResultBoolean(void* dialog_ptr);
    EXPORT void QtConsumeDialog(void* dialog_ptr, void (*cb)(uint64_t), uint64_t payload);
    EXPORT void* QtCreateInputDialog(int mode, QtVariant value, QtString title, QtString prompt);
    EXPORT void QtInputDialogUseMultilineText(void* dialog_ptr);
    EXPORT void QtInputDialogUseChoiceItems(void* dialog_ptr, QtString* items_ptr, size_t items_len);
    EXPORT QtString QtInputDialogGetTextValue(void* dialog_ptr);
    EXPORT int QtInputDialogGetIntValue(void* dialog_ptr);
    EXPORT double QtInputDialogGetDoubleValue(void* dialog_ptr);
    EXPORT void* QtCreateMessageBox(int icon, int buttons, QtString title, QtString content);
    EXPORT void QtMessageBoxSetDefaultButton(void* msgbox_ptr, int btn);
    EXPORT int QtMessageBoxGetResultButton(void* msgbox_ptr);
    EXPORT void* QtCreateFileDialog(int mode_, QtString filters);
    EXPORT int QtFileDialogGetResultFileCount(void* d_ptr);
    EXPORT QtString QtFileDialogGetResultFileItem(void* d_ptr, int index);
    EXPORT void* QtCreateLayoutRow(int spacing);
    EXPORT void* QtCreateLayoutColumn(int spacing);
    EXPORT void* QtCreateLayoutGrid(int row_spacing, int column_spacing);
    EXPORT QtGridSpan QtMakeGridSpan(int row, int column, int rowSpan, int columnSpan);
    EXPORT void QtLayoutAddLayout(void* self_ptr, void* layout_ptr, QtGridSpan span, int align);
    EXPORT void QtLayoutAddWidget(void* self_ptr, void* widget_ptr, QtGridSpan span, int align);
    EXPORT void QtLayoutAddSpacer(void* self_ptr, int width, int height, QtBool expand, QtGridSpan span, int align);
    EXPORT void QtLayoutAddLabel(void* self_ptr, QtString text, QtGridSpan span, int align);
    EXPORT void* QtCreateWidget(void* layout_ptr, int margin_x, int margin_y, int policy_x, int policy_y);
    EXPORT void* QtCreateMainWindow(void* menu_bar_ptr, void* tool_bar_ptr, void* layout_ptr, int margin_x, int margin_y, int width, int height, QtIcon icon);
    EXPORT void* QtCreateDialog(void* layout_ptr, int margin_x, int margin_y, int width, int height, QtIcon icon);
    EXPORT void* QtCreateScrollArea(int direction, void* layout_ptr, int margin_x, int margin_y);
    EXPORT void* QtCreateGroupBox(QtString title, void* layout_ptr, int margin_x, int margin_y);
    EXPORT void* QtCreateSplitter(void** widgets_ptr, size_t widgets_len);
    EXPORT void* QtCreateDynamicWidget();
    EXPORT void* QtCreateDummyFocusableWidget();
    EXPORT void QtDynamicWidgetSetWidget(void* self_ptr, void* widget_ptr);
    EXPORT void* QtCreateLabel(QtString text, int align, QtBool selectable);
    EXPORT void* QtCreateIconLabel(QtIcon icon_, int size);
    EXPORT void* QtCreateElidedLabel(QtString text);
    EXPORT void* QtCreateTextView(QtString text, int format);
    EXPORT void* QtCreateCheckBox(QtString text, QtBool checked);
    EXPORT void* QtCreateComboBox();
    EXPORT void QtComboBoxAddItem(void* b_ptr, QtIcon icon, QtString name, QtBool selected);
    EXPORT void* QtCreateComboBoxDialog(QtString title, QtString prompt);
    EXPORT void* QtComboBoxDialogGetComboBox(void* d_ptr);
    EXPORT void* QtCreatePushButton(QtIcon icon, QtString text, QtString tooltip);
    EXPORT void* QtCreateLineEdit(QtString text);
    EXPORT void* QtCreatePlainTextEdit(QtString text);
    EXPORT void* QtCreateSlider(int min, int max, int value);
    EXPORT void* QtCreateProgressBar(QtString format_, int max);
    EXPORT QtString QtClipboardReadText();
    EXPORT void QtClipboardWriteText(QtString text);
    EXPORT void* QtLwiCreateFromDefaultListWidget(size_t columns, int select_, void** headers_ptr, size_t headers_len, int stretch);
    EXPORT void* QtLwiCastToWidget(void* lwi_ptr);
    EXPORT QtString QtLwiCurrent(void* lwi_ptr, QtBool* exists);
    EXPORT QtStringList QtLwiAll(void* lwi_ptr);
    EXPORT QtStringList QtLwiSelection(void* lwi_ptr);
    EXPORT QtStringList QtLwiContiguousSelection(void* lwi_ptr);
    EXPORT void QtLwiInsert(void* lwi_ptr, int hint_, QtString pivot_, QtString key_, void** widgets_ptr, size_t widgets_len);
    EXPORT void QtLwiUpdate(void* lwi_ptr, QtString key_);
    EXPORT QtBool QtLwiMove(void* lwi_ptr, int hint_, QtString key_);
    EXPORT void QtLwiDelete(void* lwi_ptr, QtString key_);
    EXPORT void QtLwiReorder(void* lwi_ptr, QtString* order_ptr, size_t order_len);
#ifdef __cplusplus
}
#endif

#endif

