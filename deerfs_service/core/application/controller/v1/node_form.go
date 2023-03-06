package v1

// AddNodeForm 用于验证表单
type AddNodeForm struct {
	// form 标签用于 c.Bind(form)
	NodeName   string `form:"node_name" valid:"Required;MaxSize(200)"`
	UriAddress string `form:"uri_address" valid:"Required;MaxSize(200)"`
	UseCap     string `form:"use_cap" valid:"Numeric;"`
	MaxCap     string `form:"max_cap" valid:"Numeric;"`
	CreatedBy  string `form:"created_by" valid:"MaxSize(150)"` //Required;
	//DeletedOn     int    `form:"deleted_on" valid:"Range(0,1)"`
}

// EditNodeForm 用于验证表单
type EditNodeForm struct {
	//ID         int    `form:"id" valid:"Required;Min(1)"`
	NodeName   string `form:"node_name" valid:"Required;MaxSize(200)"`
	UriAddress string `form:"uri_address" valid:"Required;MaxSize(200)"`
	UseCap     string `form:"use_cap" valid:"Numeric;"`
	MaxCap     string `form:"max_cap" valid:"Numeric;"`
	ModifiedBy string `form:"modified_by" valid:"MaxSize(150)"` //Required;
}
