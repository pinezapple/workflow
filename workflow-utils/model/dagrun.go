package model

type DAGRun struct{
	DAGRunId string `gorm:"column:dagrun_id; primaryKey" json:"dagrun_id"`
	DAGId string `gorm:"column:dag_of_dagrun" json:"dag_id"`
	UserID uint32 `gorm:"column:user_of_dagrun" json:"user_id"`
	DAG DAG `gorm:"foreignKey:dag_of_dagrun;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserInfo UserInfo `gorm:"foreignKey:user_of_dagrun;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type DAG struct{
	DAGId string `gorm:"column:dag_id; primaryKey" json:"dag_id"`
	DAGName string `gorm:"column:dag_name" json:"dagname"`
	UserId uint32 `gorm:"column:user_of_dag" json:"user_id"`
	UserInfo UserInfo `gorm:"foreignKey:user_of_dag;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
