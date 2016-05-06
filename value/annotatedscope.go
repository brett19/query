package value

type annotatedScopeValue struct {
	Value
	parent AnnotatedValue
}

func NewAnnotatedScopeValue(val Value, parent AnnotatedValue) *annotatedScopeValue {
	return &annotatedScopeValue{
		Value:  val,
		parent: parent,
	}
}

func (this *annotatedScopeValue) GetValue() Value {
	return this.Value
}

func (this *annotatedScopeValue) Attachments() map[string]interface{} {
	return this.parent.Attachments()
}

func (this *annotatedScopeValue) GetAttachment(key string) interface{} {
	return this.parent.GetAttachment(key)
}

func (this *annotatedScopeValue) SetAttachment(key string, val interface{}) {
	this.parent.SetAttachment(key, val)
}

func (this *annotatedScopeValue) Covers() map[string]Value {
	return this.parent.Covers()
}

func (this *annotatedScopeValue) GetCover(key string) Value {
	return this.parent.GetCover(key)
}

func (this *annotatedScopeValue) SetCover(key string, val Value) {
	this.parent.SetCover(key, val)
}
