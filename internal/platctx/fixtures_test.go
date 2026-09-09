package platctx

const (
	subjectName    = SubjectName("payments")
	subjectNS      = Namespace("platform")
	owner          = OwnerRef("team:platform")
	operator       = OperatorRef("team:sre")
	environment    = Environment("production")
	costCenter     = CostCenter("CC-FIN-001")
	classification = DataClassification("internal")
)
