package feedback

// TODO clean errors.
var (
	ErrorUserNotFound   = UserError{"I couldn't find you in a guild with me."}
	ErrorUserNotInVoice = UserError{"I can't find you in any voice channels."}
	ErrorBotNotInVoice  = UserError{"I don't appear to be in any voice channels, you can ask me to `/join`?"}
	ErrorSoundNotFound  = UserError{"No such sound exists! Usage `/play [sound]`"}
	ErrorFatalError     = UserError{"Something unexpected has happened :("}
	ErrorAlreadyInVoice = UserError{"I am already in a voice channel for that guild."}
)
