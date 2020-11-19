package feedback

// TODO clean errors.
var (
	ErrorNoUserInVoice = UserError{"I can't find you in any voice channels."}
	ErrorNotInVoice    = UserError{"I don't appear to be in any voice channels, you can ask me to `/join`?"}
	ErrorSoundNotFound = UserError{"No such sound exists! Usage `/play [sound]`"}
)
