package bot

// This file contains all the Telegram API types used by TeleRiHa.
// These are simplified versions of the Telegram Bot API types.

// Update represents an incoming update from Telegram.
type Update struct {
	UpdateID             int64            `json:"update_id"`
	Message              *Message         `json:"message,omitempty"`
	EditedMessage        *Message         `json:"edited_message,omitempty"`
	ChannelPost          *Message         `json:"channel_post,omitempty"`
	EditedChannelPost    *Message         `json:"edited_channel_post,omitempty"`
	InlineQuery          *InlineQuery     `json:"inline_query,omitempty"`
	ChosenInlineResult   *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery        *CallbackQuery   `json:"callback_query,omitempty"`
	ShippingQuery        *ShippingQuery   `json:"shipping_query,omitempty"`
	PreCheckoutQuery     *PreCheckoutQuery `json:"pre_checkout_query,omitempty"`
	Poll                 *Poll            `json:"poll,omitempty"`
	PollAnswer           *PollAnswer      `json:"poll_answer,omitempty"`
	MyChatMember         *ChatMemberUpdated `json:"my_chat_member,omitempty"`
	ChatMember           *ChatMemberUpdated `json:"chat_member,omitempty"`
	ChatJoinRequest      *ChatJoinRequest `json:"chat_join_request,omitempty"`
}

// Message represents a Telegram message.
type Message struct {
	MessageID           int64           `json:"message_id"`
	MessageThreadID     int64           `json:"message_thread_id,omitempty"`
	From                *User           `json:"from,omitempty"`
	SenderChat          *Chat           `json:"sender_chat,omitempty"`
	Date                int64           `json:"date"`
	Chat                *Chat           `json:"chat"`
	ForwardFrom         *User           `json:"forward_from,omitempty"`
	ForwardFromChat     *Chat           `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID int64          `json:"forward_from_message_id,omitempty"`
	ForwardSignature    string          `json:"forward_signature,omitempty"`
	ForwardSenderName   string          `json:"forward_sender_name,omitempty"`
	ForwardDate         int64           `json:"forward_date,omitempty"`
	IsTopicMessage      bool            `json:"is_topic_message,omitempty"`
	IsAutomaticForward   bool            `json:"is_automatic_forward,omitempty"`
	ReplyToMessage      *Message        `json:"reply_to_message,omitempty"`
	ViaBot              *User           `json:"via_bot,omitempty"`
	EditDate            int64           `json:"edit_date,omitempty"`
	HasProtectedContent bool            `json:"has_protected_content,omitempty"`
	MediaGroupID        string          `json:"media_group_id,omitempty"`
	AuthorSignature     string          `json:"author_signature,omitempty"`
	Text                string          `json:"text,omitempty"`
	Entities            []MessageEntity `json:"entities,omitempty"`
	CaptionEntities     []MessageEntity `json:"caption_entities,omitempty"`
	Audio               *Audio          `json:"audio,omitempty"`
	Document            *Document       `json:"document,omitempty"`
	Animation           *Animation      `json:"animation,omitempty"`
	Game                *Game           `json:"game,omitempty"`
	Photo               []PhotoSize     `json:"photo,omitempty"`
	Sticker             *Sticker        `json:"sticker,omitempty"`
	Video               *Video          `json:"video,omitempty"`
	Voice               *Voice          `json:"voice,omitempty"`
	VideoNote           *VideoNote      `json:"video_note,omitempty"`
	Caption             string          `json:"caption,omitempty"`
	Contact             *Contact        `json:"contact,omitempty"`
	Location            *Location       `json:"location,omitempty"`
	Venue               *Venue          `json:"venue,omitempty"`
	Poll                *Poll           `json:"poll,omitempty"`
	Dice                *Dice           `json:"dice,omitempty"`
	NewChatMembers      []User          `json:"new_chat_members,omitempty"`
	LeftChatMember      *User           `json:"left_chat_member,omitempty"`
	NewChatTitle        string          `json:"new_chat_title,omitempty"`
	NewChatPhoto        []PhotoSize     `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto     bool            `json:"delete_chat_photo,omitempty"`
	GroupChatCreated    bool            `json:"group_chat_created,omitempty"`
	SupergroupChatCreated bool          `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated  bool            `json:"channel_chat_created,omitempty"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	MigrateToChatID     int64           `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID   int64           `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage       *Message        `json:"pinned_message,omitempty"`
	Invoice             *Invoice        `json:"invoice,omitempty"`
	SuccessfulPayment   *SuccessfulPayment `json:"successful_payment,omitempty"`
	UsersShared         *UsersShared    `json:"users_shared,omitempty"`
	ChatShared          *ChatShared     `json:"chat_shared,omitempty"`
	ConnectedWebsite    string          `json:"connected_website,omitempty"`
	WriteAccessAllowed  *WriteAccessAllowed `json:"write_access_allowed,omitempty"`
	PassportData        *PassportData   `json:"passport_data,omitempty"`
	ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"`
	BoostAdded          *ChatBoostAdded `json:"boost_added,omitempty"`
	ChatBackgroundSet   *ChatBackground `json:"chat_background_set,omitempty"`
	ForumTopicCreated    *ForumTopicCreated `json:"forum_topic_created,omitempty"`
	ForumTopicEdited     *ForumTopicEdited `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed     *ForumTopicClosed `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened   *ForumTopicReopened `json:"forum_topic_reopened,omitempty"`
	GeneralForumHide     *GeneralForumHide `json:"general_forum_hide,omitempty"`
	GeneralForumUnhide   *GeneralForumUnhide `json:"general_forum_unhide,omitempty"`
	GiveawayCreated      *Giveaway        `json:"giveaway_created,omitempty"`
	Giveaway            *Giveaway        `json:"giveaway,omitempty"`
	GiveawayWinners     *GiveawayWinners `json:"giveaway_winners,omitempty"`
	GiveawayCompleted    *GiveawayCompleted `json:"giveaway_completed,omitempty"`
	VideoChatScheduled   *VideoChatScheduled `json:"video_chat_scheduled,omitempty"`
	VideoChatStarted     *VideoChatStarted `json:"video_chat_started,omitempty"`
	VideoChatEnded       *VideoChatEnded   `json:"video_chat_ended,omitempty"`
	VideoChatParticipantsInvited *VideoChatParticipantsInvited `json:"video_chat_participants_invited,omitempty"`
	WebAppData          *WebAppData      `json:"web_app_data,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// User represents a Telegram user.
type User struct {
	ID                          int64  `json:"id"`
	IsBot                      bool   `json:"is_bot"`
	FirstName                   string `json:"first_name"`
	LastName                    string `json:"last_name,omitempty"`
	Username                    string `json:"username,omitempty"`
	LanguageCode                string `json:"language_code,omitempty"`
	IsPremium                   bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu       bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups               bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages     bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries       bool   `json:"supports_inline_queries,omitempty"`
}

// Chat represents a Telegram chat.
type Chat struct {
	ID                          int64  `json:"id"`
	Type                        string `json:"type"`
	Title                       string `json:"title,omitempty"`
	Username                    string `json:"username,omitempty"`
	FirstName                   string `json:"first_name,omitempty"`
	LastName                    string `json:"last_name,omitempty"`
	IsForum                     bool   `json:"is_forum,omitempty"`
	Photo                       *ChatPhoto `json:"photo,omitempty"`
	ActiveUsernames             []string `json:"active_usernames,omitempty"`
	EmojiStatusCustomEmojiID    string `json:"emoji_status_custom_emoji_id,omitempty"`
	Bio                         string `json:"bio,omitempty"`
	HasPrivateForwards          bool   `json:"has_private_forwards,omitempty"`
	HasRestrictedVoiceAndVideoMessages bool `json:"has_restricted_voice_and_video_messages,omitempty"`
	JoinToSendMessages          bool   `json:"join_to_send_messages,omitempty"`
	JoinByRequest               bool   `json:"join_by_request,omitempty"`
	HasAggressiveAntiSpamEnabled bool   `json:"has_aggressive_anti_spam_enabled,omitempty"`
	HideNonPublicMembers        bool   `json:"hide_non_public_members,omitempty"`
	AllMembersAreAdministrators  bool   `json:"all_members_are_administrators,omitempty"`
	CanSetStickerSet           bool   `json:"can_set_sticker_set,omitempty"`
	CustomEmojiStickerSetName   string `json:"custom_emoji_sticker_set_name,omitempty"`
	LinkedChatID                int64  `json:"linked_chat_id,omitempty"`
	Location                    *ChatLocation `json:"location,omitempty"`
}

// MessageEntity represents a part of a message text.
type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	URL           string `json:"url,omitempty"`
	User          *User  `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// PhotoSize represents a photo in different sizes.
type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Audio represents an audio file.
type Audio struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	Performer    string `json:"performer,omitempty"`
	Title        string `json:"title,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
}

// Document represents a document file.
type Document struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Video represents a video file.
type Video struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

// Sticker represents a sticker.
type Sticker struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Type         string `json:"type"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	IsAnimated   bool   `json:"is_animated,omitempty"`
	IsVideo      bool   `json:"is_video,omitempty"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	Emoji        string `json:"emoji,omitempty"`
	SetName      string `json:"set_name,omitempty"`
	PremiumAnimation *File `json:"premium_animation,omitempty"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
	NeedsRepainting bool `json:"needs_repainting,omitempty"`
	IsFullBody  bool   `json:"is_full_body,omitempty"`
}

// InlineQuery represents an inline query.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}

// CallbackQuery represents a callback query from inline buttons.
type CallbackQuery struct {
	ID              string       `json:"id"`
	From            *User        `json:"from"`
	Message         *Message     `json:"message,omitempty"`
	InlineMessageID string       `json:"inline_message_id,omitempty"`
	ChatInstance    string       `json:"chat_instance"`
	Data            string       `json:"data,omitempty"`
	GameShortName   string       `json:"game_short_name,omitempty"`
}

// ChosenInlineResult represents a chosen inline result.
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            *User     `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageID string    `json:"inline_message_id,omitempty"`
	Query           string    `json:"query"`
}

// ChatPhoto represents a chat photo.
type ChatPhoto struct {
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`
	BigFileID         string `json:"big_file_id"`
	BigFileUniqueID   string `json:"big_file_unique_id"`
}

// ChatMemberUpdated represents a chat member update.
type ChatMemberUpdated struct {
	Chat     *Chat       `json:"chat"`
	From     *User       `json:"from"`
	Date     int64       `json:"date"`
	OldChatMember *ChatMember `json:"old_chat_member"`
	NewChatMember *ChatMember `json:"new_chat_member"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
	ViaChatFolderInviteLink bool `json:"via_chat_folder_invite_link,omitempty"`
}

// ChatMember represents a chat member.
type ChatMember struct {
	Status   string `json:"status"`
	User     *User  `json:"user"`
	CustomTitle string `json:"custom_title,omitempty"`
	UntilDate int64  `json:"until_date,omitempty"`
	CanBeEdited bool  `json:"can_be_edited,omitempty"`
	CanChangeInfo bool `json:"can_change_info,omitempty"`
	CanPostMessages bool `json:"can_post_messages,omitempty"`
	CanEditMessages bool `json:"can_edit_messages,omitempty"`
	CanDeleteMessages bool `json:"can_delete_messages,omitempty"`
	CanRestrictMembers bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers bool `json:"can_promote_members,omitempty"`
	CanManageChat bool `json:"can_manage_chat,omitempty"`
	CanManageVideoChats bool `json:"can_manage_video_chats,omitempty"`
	IsAnonymous bool `json:"is_anonymous,omitempty"`
	CanManageTopics bool `json:"can_manage_topics,omitempty"`
}

// Location represents a geographic location.
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	HorizontalAccuracy float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod int `json:"live_period,omitempty"`
	Heading    int `json:"heading,omitempty"`
	ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
}

// InlineKeyboardMarkup represents an inline keyboard.
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represents a button in an inline keyboard.
type InlineKeyboardButton struct {
	Text            string `json:"text"`
	URL             string `json:"url,omitempty"`
	LoginURL        *LoginURL `json:"login_url,omitempty"`
	CallbackData    string `json:"callback_data,omitempty"`
	WebApp          *WebAppInfo `json:"web_app,omitempty"`
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
	CallbackGame    *CallbackGame `json:"callback_game,omitempty"`
	Pay             bool `json:"pay,omitempty"`
}

// ReplyKeyboardMarkup represents a reply keyboard.
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool              `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool              `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string      `json:"input_field_placeholder,omitempty"`
	Selective       bool              `json:"selective,omitempty"`
}

// KeyboardButton represents a button in a reply keyboard.
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo `json:"web_app,omitempty"`
}

// File represents a file from Telegram.
type File struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int64  `json:"file_size,omitempty"`
	FilePath     string `json:"file_path,omitempty"`
}

// Additional types for completeness

type Animation struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type Voice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type VideoNote struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Length       int    `json:"length"`
	Duration     int    `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileSize     int64  `json:"file_size,omitempty"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserID      int64  `json:"user_id,omitempty"`
	Vcard       string `json:"vcard,omitempty"`
}

type Venue struct {
	Location        *Location `json:"location"`
	Title           string    `json:"title"`
	Address         string    `json:"address"`
	FoursquareID    string    `json:"foursquare_id,omitempty"`
	FoursquareType   string    `json:"foursquare_type,omitempty"`
	GooglePlaceID   string    `json:"google_place_id,omitempty"`
	GooglePlaceType string    `json:"google_place_type,omitempty"`
}

type Poll struct {
	ID                    string       `json:"id"`
	Question              string       `json:"question"`
	Options               []PollOption `json:"options"`
	TotalVoterCount       int          `json:"total_voter_count"`
	IsClosed              bool         `json:"is_closed"`
	IsAnonymous           bool         `json:"is_anonymous"`
	Type                  string       `json:"type"`
	AllowsMultipleAnswers bool         `json:"allows_multiple_answers"`
	CorrectOptionID       int          `json:"correct_option_id,omitempty"`
	Explanation           string       `json:"explanation,omitempty"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            int          `json:"open_period"`
	CloseDate             int64        `json:"close_date,omitempty"`
}

type PollOption struct {
	Text            string `json:"text"`
	VoterCount      int    `json:"voter_count"`
	VotePercentage  int    `json:"vote_percentage,omitempty"`
}

type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int64  `json:"total_amount"`
}

type SuccessfulPayment struct {
	Currency                string `json:"currency"`
	TotalAmount             int64  `json:"total_amount"`
	InvoicePayload          string `json:"invoice_payload"`
	ShippingOptionID        string `json:"shipping_option_id,omitempty"`
	OrderInfo               *OrderInfo `json:"order_info,omitempty"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
}

type ShippingQuery struct {
	ID              string `json:"id"`
	From            *User  `json:"from"`
	InvoicePayload  string `json:"invoice_payload"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

type PreCheckoutQuery struct {
	ID               string `json:"id"`
	From             *User  `json:"from"`
	Currency         string `json:"currency"`
	TotalAmount      int64  `json:"total_amount"`
	InvoicePayload   string `json:"invoice_payload"`
	ShippingOptionID string `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

type ChatJoinRequest struct {
	Chat     *Chat `json:"chat"`
	From     *User `json:"from"`
	Date     int64 `json:"date"`
	Bio      string `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

type ChatBoostAdded struct {
	BoostID string `json:"boost_id"`
}

type WebAppData struct {
	Data string `json:"data"`
	ButtonText string `json:"button_text"`
}

// Simplified types for common use cases

type ChatInviteLink struct {
	Link string `json:"link"`
}

type OrderInfo struct {
	Name            string `json:"name,omitempty"`
	PhoneNumber     string `json:"phone_number,omitempty"`
	Email           string `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

type UsersShared struct {
	RequestID int64 `json:"request_id"`
	UserIDs   []int64 `json:"user_ids"`
}

type ChatShared struct {
	RequestID int64 `json:"request_id"`
	ChatID    int64 `json:"chat_id"`
}

type WriteAccessAllowed struct {
	FromRequest bool `json:"from_request"`
	WebAppName string `json:"web_app_name,omitempty"`
	FromAttachmentMenu bool `json:"from_attachment_menu,omitempty"`
}

type PassportData struct {
	Credentials *EncryptedCredentials `json:"credentials"`
	Data        []EncryptedPassportElement `json:"data"`
}

type EncryptedCredentials struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	Secret string `json:"secret"`
}

type EncryptedPassportElement struct {
	Type string `json:"type"`
	Hash string `json:"hash"`
	Data string `json:"data,omitempty"`
}

type ProximityAlertTriggered struct {
	Traveler *User `json:"traveler"`
	Watcher  *User `json:"watcher"`
	Distance int    `json:"distance"`
}

type ChatBackground struct {
	Type string `json:"type"`
}

type ForumTopicCreated struct {
	Name string `json:"name"`
	IconColor int `json:"icon_color"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

type ForumTopicEdited struct {
	Name string `json:"name,omitempty"`
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

type ForumTopicClosed struct {
}

type ForumTopicReopened struct {
}

type GeneralForumHide struct {
}

type GeneralForumUnhide struct {
}

type Giveaway struct {
	Chats          []Chat `json:"chats"`
	Dates          *GiveawayDates `json:"dates"`
	WinnersSelectionDate int64 `json:"winners_selection_date"`
	WinnerCount    int    `json:"winner_count"`
	OnlyNewMembers bool   `json:"only_new_members,omitempty"`
	HasPublicWinners bool `json:"has_public_winners,omitempty"`
	PrizeDescription string `json:"prize_description,omitempty"`
	CountryCodes   []string `json:"country_codes,omitempty"`
	PremiumSubscriptionMonthCount []int `json:"premium_subscription_month_count,omitempty"`
}

type GiveawayDates struct {
	StartDate int64 `json:"start_date"`
	EndDate   int64 `json:"end_date"`
}

type GiveawayWinners struct {
	Chat *Chat `json:"chat"`
	Winners []User `json:"winners"`
	WinnerCount int `json:"winner_count"`
	UnclaimedPrizeCount int `json:"unclaimed_prize_count,omitempty"`
	GiveawayMessage *Message `json:"giveaway_message,omitempty"`
}

type GiveawayCompleted struct {
	WinnerCount int `json:"winner_count"`
	UnclaimedPrizeCount int `json:"unclaimed_prize_count,omitempty"`
	GiveawayMessage *Message `json:"giveaway_message"`
	Winners []User `json:"winners"`
	AdditionalChatCount int `json:"additional_chat_count,omitempty"`
	PremiumSubscriptionMonthCount int `json:"premium_subscription_month_count,omitempty"`
}

type VideoChatScheduled struct {
	StartDate int64 `json:"start_date"`
}

type VideoChatStarted struct {
}

type VideoChatEnded struct {
	Duration int `json:"duration"`
}

type VideoChatParticipantsInvited struct {
	Participants []User `json:"participants"`
}

type PollAnswer struct {
	PollID string `json:"poll_id"`
	User   *User  `json:"user"`
	OptionIDs []int `json:"option_ids"`
}

type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

type LoginURL struct {
	URL        string `json:"url"`
	ForwardText string `json:"forward_text,omitempty"`
	BotUsername string `json:"bot_username,omitempty"`
	RequestWriteAccess bool `json:"request_write_access,omitempty"`
}

type CallbackGame struct {
}

type WebAppInfo struct {
	URL string `json:"url"`
}

type KeyboardButtonPollType struct {
	Type string `json:"type"`
}

type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float64 `json:"x_shift"`
	YShift float64 `json:"y_shift"`
	Scale  float64 `json:"scale"`
}

type ChatLocation struct {
	Location *Location `json:"location"`
	Address  string `json:"address"`
}

type Game struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Photo        []PhotoSize `json:"photo,omitempty"`
	Text         string `json:"text,omitempty"`
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
	Animation    *Animation `json:"animation,omitempty"`
}

type VoiceChatScheduled struct {
	StartDate int64 `json:"start_date"`
}

type VoiceChatStarted struct {
}

type VoiceChatEnded struct {
	Duration int `json:"duration"`
}

type VoiceChatParticipantsInvited struct {
	Participants []User `json:"participants"`
}
