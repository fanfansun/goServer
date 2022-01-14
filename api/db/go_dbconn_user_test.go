package db

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tempvid string

func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

func TestMain(m *testing.M) {
	clearTables()
	m.Run()
	clearTables()
}

func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", TestAddUserInformation)
	t.Run("Get", TestGetUserInformation)
	t.Run("Del", TestDeleteUser)

}

func TestAddUserInformation(t *testing.T) {
	err := AddUserInformation("avenssi", "123")
	if err != nil {
		t.Errorf("ERROR of AddUser: %v", err)
	}
}

func TestGetUserInformation(t *testing.T) {
	pwd, err := GetUserInformation("avenssi")
	if err != nil {
		t.Errorf("ERROR of GetUser: %v", err)
	}
	fmt.Printf("pwd is %s\n", pwd)

}

func TestDeleteUser(t *testing.T) {
	err := DeleteUser("avenssi", "123")
	if err != nil {
		t.Errorf("Error of DelUser: %v", err)
	}

}

func TestRegetUser(t *testing.T) {
	pwd, err := GetUserInformation("avenssi")
	if err != nil {
		t.Errorf("ERROR of GetUser: %v", err)
	}
	fmt.Printf("pwd is %v\n", pwd)
}

func TestDBconnUser(t *testing.T) {
	fmt.Println("123")
}

func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", TestAddUserInformation)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

func testAddVideoInfo(t *testing.T) {
	vi, err := AddvideoInformation(1, "my-video")
	if err != nil {
		t.Errorf("Error of AddVideoInfo: %v", err)
	}
	tempvid = vi.Id
}

func testGetVideoInfo(t *testing.T) {
	_, err := GetvideoInformation(tempvid)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

func testDeleteVideoInfo(t *testing.T) {
	err := Deletevideo(tempvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoInfo: %v", err)
	}
}

func testRegetVideoInfo(t *testing.T) {
	vi, err := GetvideoInformation(tempvid)
	if err != nil || vi != nil {
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", TestAddUserInformation)
	t.Run("AddCommnets", testAddComments)
	t.Run("ListComments", testListComments)
}

func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"

	err := AddcommentInformation(vid, aid, content)

	if err != nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

func testListComments(t *testing.T) {
	vid := "12345"
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))
	fmt.Print("time now UnixNano is \n", time.Now().UnixNano())
	// 1642120435095357200 -> 1642120435
	fmt.Print("to is \n", to)
	res, err := ListComments(vid, from, to)
	print("res is \n", res)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}

	for i, ele := range res {
		fmt.Printf("comment: %d, %v \n", i, ele)
	}
}
