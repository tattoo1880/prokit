package codec

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
	"time"

	pb "github.com/tattoo1880/protkit/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// 构造一个完整的 User 用于测试
func newTestUser() *pb.User {
	now := timestamppb.New(time.Date(2026, 4, 13, 12, 0, 0, 0, time.UTC))
	return &pb.User{
		Id:       42,
		Username: "alice",
		Email:    "alice@example.com",
		Phone:    "13800138000",
		Gender:   pb.Gender_GENDER_FEMALE,
		Age:      28,
		Address: &pb.Address{
			Province: "浙江省",
			City:     "杭州市",
			District: "西湖区",
			Street:   "文三路 123 号",
			ZipCode:  "310000",
		},
		Tags:      []string{"developer", "gopher"},
		Extra:     map[string]string{"team": "backend", "level": "senior"},
		Status:    pb.UserStatus_USER_STATUS_ACTIVE,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ---------- Marshal / Unmarshal ----------

func TestMarshalUnmarshal(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 4, 13, 12, 0, 0, 0, time.UTC))
	user := &pb.User{
		Id:       1,
		Username: "marshal_test",
		Email:    "marshal@example.com",
		Phone:    "13900000001",
		Gender:   pb.Gender_GENDER_MALE,
		Age:      25,
		Address: &pb.Address{
			Province: "北京市",
			City:     "北京市",
			District: "海淀区",
			Street:   "中关村大街 1 号",
			ZipCode:  "100080",
		},
		Tags:      []string{"test", "marshal"},
		Extra:     map[string]string{"source": "unit_test"},
		Status:    pb.UserStatus_USER_STATUS_ACTIVE,
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, err := Marshal(user)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Marshal returned empty bytes")
	}

	got := &pb.User{}
	if err := Unmarshal(data, got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !proto.Equal(user, got) {
		t.Errorf("round-trip mismatch:\n  want: %v\n  got:  %v", user, got)
	}
}

func TestMarshalNil(t *testing.T) {
	// var user *pb.User 是 nil 指针，对其字段赋值（user.Age = 30）会触发 nil pointer panic，
	// 这里直接验证 Marshal(nil) 应返回 error
	var user *pb.User
	data, err := Marshal(user)
	if err == nil {
		t.Fatal("expected error for nil message")
	}
	fmt.Println(data)
}

func TestUnmarshalNil(t *testing.T) {
	if err := Unmarshal([]byte{}, nil); err == nil {
		t.Fatal("expected error for nil message")
	}
}

// ---------- 各字段类型验证 ----------

func TestFieldTypes(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	user := &pb.User{
		Id:       100,
		Username: "field_types_test",
		Email:    "fields@example.com",
		Phone:    "13700000002",
		Gender:   pb.Gender_GENDER_FEMALE,
		Age:      30,
		Address: &pb.Address{
			Province: "上海市",
			City:     "上海市",
			District: "浦东新区",
			Street:   "张江高科 88 号",
			ZipCode:  "201203",
		},
		Tags:      []string{"engineer", "reader"},
		Extra:     map[string]string{"dept": "infra", "rank": "p7"},
		Status:    pb.UserStatus_USER_STATUS_INACTIVE,
		CreatedAt: now,
		UpdatedAt: now,
	}

	data, err := Marshal(user)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	got := &pb.User{}
	if err := Unmarshal(data, got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	// 标量
	if got.Id != 100 {
		t.Errorf("Id = %d, want 100", got.Id)
	}
	if got.Username != "field_types_test" {
		t.Errorf("Username = %q, want field_types_test", got.Username)
	}
	if got.Age != 30 {
		t.Errorf("Age = %d, want 30", got.Age)
	}

	// 枚举
	if got.Gender != pb.Gender_GENDER_FEMALE {
		t.Errorf("Gender = %v, want FEMALE", got.Gender)
	}
	if got.Status != pb.UserStatus_USER_STATUS_INACTIVE {
		t.Errorf("Status = %v, want INACTIVE", got.Status)
	}

	// 嵌套消息
	if got.Address == nil || got.Address.City != "上海市" {
		t.Errorf("Address.City = %v, want 上海市", got.GetAddress().GetCity())
	}

	// repeated
	if len(got.Tags) != 2 || got.Tags[0] != "engineer" {
		t.Errorf("Tags = %v, want [engineer reader]", got.Tags)
	}

	// map
	if got.Extra["dept"] != "infra" {
		t.Errorf("Extra[dept] = %q, want infra", got.Extra["dept"])
	}

	// Timestamp
	if got.CreatedAt.AsTime().Year() != 2026 {
		t.Errorf("CreatedAt year = %d, want 2026", got.CreatedAt.AsTime().Year())
	}
}

// ---------- 默认值 / 零值 ----------

func TestEmptyUser(t *testing.T) {
	user := &pb.User{}

	data, err := Marshal(user)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	got := &pb.User{}
	if err := Unmarshal(data, got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if got.Id != 0 || got.Username != "" || got.Gender != pb.Gender_GENDER_UNSPECIFIED {
		t.Errorf("zero-value user mismatch: %v", got)
	}
}

// ---------- UserList ----------

func TestUserList(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC))
	list := &pb.UserList{
		Users: []*pb.User{
			{
				Id:       200,
				Username: "list_user_a",
				Email:    "a@example.com",
				Phone:    "13600000003",
				Gender:   pb.Gender_GENDER_MALE,
				Age:      22,
				Address: &pb.Address{
					Province: "广东省",
					City:     "深圳市",
					District: "南山区",
					Street:   "科技园路 10 号",
					ZipCode:  "518000",
				},
				Tags:      []string{"new"},
				Extra:     map[string]string{"role": "intern"},
				Status:    pb.UserStatus_USER_STATUS_ACTIVE,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{Id: 201, Username: "list_user_b"},
		},
		TotalCount: 2,
	}

	data, err := Marshal(list)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	got := &pb.UserList{}
	if err := Unmarshal(data, got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if !proto.Equal(list, got) {
		t.Errorf("UserList round-trip mismatch")
	}
	if got.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", got.TotalCount)
	}
	if got.Users[1].Username != "list_user_b" {
		t.Errorf("Users[1].Username = %q, want list_user_b", got.Users[1].Username)
	}
}

// ---------- Frame 帧协议 ----------

func TestWriteReadFrame(t *testing.T) {
	now := timestamppb.New(time.Date(2026, 3, 15, 10, 30, 0, 0, time.UTC))
	user := &pb.User{
		Id:       300,
		Username: "frame_test",
		Email:    "frame@example.com",
		Phone:    "13500000004",
		Gender:   pb.Gender_GENDER_FEMALE,
		Age:      35,
		Address: &pb.Address{
			Province: "四川省",
			City:     "成都市",
			District: "武侯区",
			Street:   "天府大道 200 号",
			ZipCode:  "610041",
		},
		Tags:      []string{"frame", "tcp"},
		Extra:     map[string]string{"protocol": "v2"},
		Status:    pb.UserStatus_USER_STATUS_ACTIVE,
		CreatedAt: now,
		UpdatedAt: now,
	}

	var buf bytes.Buffer
	if err := WriteFrame(&buf, user); err != nil {
		t.Fatalf("WriteFrame: %v", err)
	}

	got := &pb.User{}
	if err := ReadFrame(bufio.NewReader(&buf), got); err != nil {
		t.Fatalf("ReadFrame: %v", err)
	}

	if !proto.Equal(user, got) {
		t.Errorf("frame round-trip mismatch:\n  want: %v\n  got:  %v", user, got)
	}
}

func TestMultipleFrames(t *testing.T) {
	users := []*pb.User{
		{Id: 1, Username: "alice"},
		{Id: 2, Username: "bob"},
		{Id: 3, Username: "charlie"},
	}

	var buf bytes.Buffer
	for _, u := range users {
		if err := WriteFrame(&buf, u); err != nil {
			t.Fatalf("WriteFrame(%s): %v", u.Username, err)
		}
	}

	reader := bufio.NewReader(&buf)
	for i, want := range users {
		got := &pb.User{}
		if err := ReadFrame(reader, got); err != nil {
			t.Fatalf("ReadFrame[%d]: %v", i, err)
		}
		if !proto.Equal(want, got) {
			t.Errorf("frame[%d] mismatch: want %v, got %v", i, want, got)
		}
	}
}