package types

const (
	EventTypeAddJob            = "add_job"
	EventTypeRemoveJob         = "remove_job"
	EventTypeChangeJobCID      = "change_job_cid"
	EventTypeChangeJobLabel    = "change_job_label"
	EventTypeChangeJobCallData = "change_job_calldata"
	EventTypeChangeJobGasPrice = "change_job_gasprice"
	EventTypeChangeJobPeriod   = "change_job_period"
	EventTypeChangeJobBlock    = "change_job_block"

	AttributeKeyJobContract = "contract"
	AttributeKeyJobTrigger  = "trigger"
	AttributeKeyJobLoad     = "load"
	AttributeKeyJobLabel    = "label"
	AttributeKeyJobCID      = "cid"
	AttributeKeyJobCallData = "call_data"
	AttributeKeyJobGasPrice = "gas_price"
	AttributeKeyJobPeriod   = "period"
	AttributeKeyJobBlock    = "block"

	AttributeValueCategory   		= ModuleName
)
