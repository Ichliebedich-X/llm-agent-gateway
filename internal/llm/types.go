package llm

// MessageRole 表示对话消息的角色
type MessageRole string

const (
	RoleSystem    MessageRole = "system"    //系统规则
	RoleUser      MessageRole = "user"      //用户输入
	RoleAssistant MessageRole = "assistant" //模型回答
	RoleTool      MessageRole = "tool"      //工具执行结果
)

// Message 表示一条对话消息
type Message struct {
	Role    MessageRole //定义的四种角色
	Content string      //消息正文
}

// ChatRequest 表示一次模型调用请求
type ChatRequest struct {
	Model     string    //模型名称
	Message   []Message //一次请求包含多条消息，[]Message 就是一个 Message 类型的切片
	MaxTokens int
}

// Usage 表示模型调用的Token用量
type Usage struct {
	InputTokens  int
	OutputTokens int
}

// ChatResponse 表示模型调用的统一响应
type ChatResponse struct {
	Content  string
	Provider string
	Model    string //实际使用了哪个模型
	Usage    Usage  //Usage结构体：Token 用量本身就是一组相关数据。独立出来后，响应结构会更清晰，未来增加其他用量字段也更方便
}
