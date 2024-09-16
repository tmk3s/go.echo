package handler

// handlerのリスト。新規にhandlerを追加するたびにこちらにも追加
type AppHandler struct {
	AuthHandler AuthHandler
	TodoHandler TodoHandler
	UserHandler UserHandler
}
