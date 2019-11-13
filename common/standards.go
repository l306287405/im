package common

const(

	// ----- 底层错误区块 Begin -----
	CONTACTS_BROKEN="1001" //长连接错误码
	ONLINE_FAILD="1002"
	OFFLINE_FAILD="1003"
	// ----- 底层错误区块 End -----

	// ----- 校验错误区块 Begin -----
	ATTRIBUTION_ERROR="2001" //归属错误

	// ----- 校验错误区块 End -----

	// ----- 业务逻辑错误区块 Begin -----
	SQL_ERROR="3000"
	SQL_INSERT_FAILD="3001"
	SQL_UPDATE_FAILD="3002"

	// ----- 业务逻辑错误区块 End -----
)