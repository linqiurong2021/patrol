package model

import (
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Model struct {
	engine *xorm.Engine
}
//
type User struct {
	ID       int64 `json:"id" xorm:"id pk autoincr"`
}

// 巡河表
type Patrol struct {
	ID       int64 `json:"id" xorm:"id pk autoincr"` // 不要写 xorm:"id" 否则插入会提示  cannot insert into column "id"
	RequestID string `json:"request_id" xorm:"request_id"`
	UserID 	 int64  `json:"user_id" xorm:"user_id"`
	Memo 	string  `json:"memo" xorm:"memo"`
	Polyline string `json:"polyline" xorm:"polyline"`
	CreateAt time.Time `json:"create_at" xorm:"created"`
	UpdateAt time.Time `json:"update_at" xorm:"updated"`
	DeleteAt time.Time `json:"delete_at" xorm:"deleted"`
}

// 获取用户信息
func (m *Model) GetPatrol(UserID int64,ID int64) ( *Patrol,  error) {
	var patrol Patrol;
	//_, err := m.engine.Table(&Patrol{}).Where("user_id = ?",UserID).Where("id = ?", ID).Cols("id,user_id,memo,st_astext(line),create_at,update_at").Where("delete_at is null").Get(&patrol)
	_, err := m.engine.Table(&Patrol{}).Where("id = ?", ID).Where("delete_at is null").Select("id,user_id,memo,st_astext(polyline) as polyline,create_at,update_at").Get(&patrol);
	if err != nil {
		return nil, err
	}
	log.Printf("patrol: %#v\n",patrol)
	return &patrol, nil
}

// 获取用户信息
func (m *Model) GetPatrolList(UserID int64,Page int, PageSize int) ( []*Patrol,  error) {
	var patrol []*Patrol;
	start := (Page-1) * PageSize
	log.Printf("start %d, pageSize: %d \n", start,PageSize)

	err := m.engine.Table(&Patrol{}).Where("user_id = ?", UserID).Cols("id,user_id,memo,create_at,update_at").Where("delete_at is null").Limit(PageSize,start).Find(&patrol)
	if err != nil {
		return nil, err
	}
	//log.Printf("数据库中查询到%d行\n",len(patrol))
	return patrol, nil
}

// 获取用户巡河总条数
func (m *Model) GetPatrolCount(UserID int64) ( int64,  error) {
	var count int64
	count,err := m.engine.Table(&Patrol{}).Where("user_id = ?", UserID).Cols("id").Where("delete_at is null").Count()
	if err != nil {
		return 0, err
	}
	//log.Printf("数据库中查询到%d行\n",len(patrol))
	return count, nil
}




// 获取用户信息
func (m *Model) DelPatrol(ID int64) ( int64,  error) {

	return m.engine.Table(&Patrol{}).Where("id = ?", ID).Delete(&Patrol{})
}



// 创建巡河
func (m *Model) PostPatrol(patrolData *Patrol) (outPatrol *Patrol, err error) {
	log.Printf("insert data %#v\n",patrolData)
	result, err := m.engine.Exec("INSERT INTO Patrol(request_id,user_id,memo,polyline,create_at) values(?,?,?,st_geomfromtext(?),?) RETURNING ID",patrolData.RequestID,patrolData.UserID,patrolData.Memo,patrolData.Polyline,time.Now())
	//
	if err != nil {
		log.Printf("insert data error %s\n", err)
		return nil, err
	}
	//
	count , err := result.RowsAffected()
	if err != nil {
		log.Printf("insert data result.RowsAffected  %s\n", err)
		return nil, err
	}
	log.Printf("insert data result.count  %s\n", err)

	if count > 0 {
		patrolData.ID = count
		return patrolData, nil
	}
	return nil, nil
}


// 更新用户
func (m *Model) UpdatePatrol (patrol *Patrol) (affected int64, err error) {
	log.Printf("update user data %#v\n",patrol)
	return m.engine.Where("id = ?",patrol.ID).AllCols().Omit("id").Update(patrol)
}

//
func NewModel(dataSourceName string) *Model {
	// 连接字符串
	//dataSourceName := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Config.Postgre.Host, conf.Config.Postgre.Port, conf.Config.Postgre.User, conf.Config.Postgre.Password, conf.Config.Postgre.DBName) // db
	log.Printf("### dataSourceName:%s\n",dataSourceName)
	engine, err := xorm.NewEngine("postgres", dataSourceName)
	if err != nil {
		//
		log.Fatalf("init db error: %s", err)
	}
	//
	//err = engine.Sync2(new(Patrol))
	//if err != nil {
	//	log.Fatalf("engine sync error %s\n",err)
	//}
	return &Model{
		engine: engine,
	}
}
