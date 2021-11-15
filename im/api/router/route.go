package route

import "go_im/im/message"

type ResponseFunc = func(uid int64, message2 *message.Message)
