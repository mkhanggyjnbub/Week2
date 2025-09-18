package main

import (
	"baitapweek2/Db"
	"baitapweek2/Middleware"
	"baitapweek2/Query"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	Db.Connect()

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("02/01/2006") // dd/mm/yyyy để hiển thị
		},
		"formatDateInput": func(t time.Time) string {
			return t.Format("2006-01-02") // yyyy-mm-dd để dùng trong <input type="date">
		},
	})

	// Nạp template HTML (toàn bộ file trong thư mục templates)
	r.LoadHTMLGlob("templates/*")
	// Trang chủ: render index.html
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	// hiển thị trang login
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	//hiển thị trang home
	r.GET("/home", Middleware.AuthMiddleware("user"), func(c *gin.Context) {
		category, err := Query.GetCategory()
		if err != nil {
			c.String(http.StatusInternalServerError, "Lỗi lấy category")
			return
		}

		tasks, err := Query.GetTasks()
		if err != nil {
			c.String(http.StatusInternalServerError, "Lỗi lấy Tasks")
			return
		}

		c.HTML(http.StatusOK, "home.html", gin.H{
			"Category": category,
			"Tasks":    tasks,
		})

	})

	// đăng ký
	r.POST("/Register", func(c *gin.Context) {
		user := Db.Users{
			UserName:     c.PostForm("username"),
			Email:        c.PostForm("email"),
			PasswordHash: c.PostForm("password"),
			Role:         "user",
		}
		count, err := Query.InsertUser(user)

		if err != nil || count == 0 {
			c.JSON(http.StatusUnauthorized, "error, Đăng ký thất bại")
			return
		}

		c.Redirect(http.StatusFound, "/login")
	})

	// đăng nhập
	r.POST("/login", func(c *gin.Context) {
		email := c.PostForm("email")
		pas := c.PostForm("password")
		role := c.PostForm("role")
		fmt.Println(email)
		fmt.Println(pas)

		// kiểm tra người dùng đăng nhập hợp lệ ko
		UserQ, err := Query.CheckLogin(email, pas, role)

		if err != nil {
			// c.JSON(http.StatusUnauthorized, "error")
			c.Redirect(http.StatusFound, "/login")
			return

		}
		// nếu hợp lệ thì trả về 1 JWT
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Middleware.Claims{
			UserID: uint(UserQ.UserID),
			Email:  UserQ.Email,
			Role:   UserQ.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(Middleware.JwtKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
			return
		}
		// lưu vào cookie
		c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)

		// c.JSON(http.StatusOK, gin.H{"token": tokenString})
		fmt.Println(tokenString)
		if role == "user" {
			c.Redirect(http.StatusFound, "/home")
		} else {

			c.JSON(http.StatusOK, "chào mừng bạn đến với trang admin")
		}
	})

	// insert task and
	r.POST("/createTask", Middleware.AuthMiddleware("user"), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		userIdPase := userID.(uint)
		parseDueDate, err := time.Parse("2006-01-02", c.PostForm("dueDate"))

		if err != nil {
			return

		}
		parseCategory, err := strconv.Atoi(c.PostForm("categoryID"))

		if err != nil {
			return

		}
		task := Db.Task{
			Title:       c.PostForm("title"),
			Description: c.PostForm("description"),
			DueDate:     parseDueDate,
			CategoryID:  parseCategory,
			Status:      "pending",
			UserID:      int(userIdPase),
		}

		result, err := Query.InsertTask(task)
		if result == 0 {
			fmt.Println("insert task error")
		}
		c.Redirect(http.StatusFound, "/home")
	})

	// In ra giao diện update Task
	r.GET("/editTask", Middleware.AuthMiddleware("user"), func(c *gin.Context) {
		parseId, err := strconv.Atoi(c.Query("id"))
		task, err := Query.GetTasksById(parseId)
		category, err := Query.GetCategory()
		if err != nil {
			return
		}

		c.HTML(200, "updateTask.html", gin.H{
			"Stask":    task,
			"Category": category,
		})
	})

	// cập nhật task
	r.POST("/editTask", Middleware.AuthMiddleware("user"), func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.PostForm("task_id"))
		if err != nil {
			ctx.JSON(400, "Task ID không hợp lệ")
			return
		}

		parseDueDate, err := time.Parse("2006-01-02", ctx.PostForm("due_date"))
		if err != nil {
			ctx.JSON(400, "Ngày không hợp lệ")
			return
		}

		parseCategory, err := strconv.Atoi(ctx.PostForm("category_id"))
		if err != nil {
			ctx.JSON(400, "Category ID không hợp lệ")
			return
		}

		task := Db.Task{
			TaskID:      id,
			Title:       ctx.PostForm("title"),
			Description: ctx.PostForm("description"),
			CategoryID:  parseCategory,
			DueDate:     parseDueDate,
			Status:      ctx.PostForm("status"),
		}

		result, err := Query.EditTask(task)
		if err != nil {
			ctx.JSON(500, "Cập nhật thất bại")
			return
		}
		if result == 0 {
			ctx.JSON(404, "Task không tồn tại")
			return
		}

		ctx.Redirect(http.StatusFound, "/home")
	})

	r.GET("/deleteTask", Middleware.AuthMiddleware("user"), func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Query("id"))

		if err != nil {
			return
		}

		result, err := Query.DeleteTask(id)

		if result == 0 {
			ctx.JSON(500, gin.H{"error": "khong the xoa"})

		}
		ctx.Redirect(http.StatusFound, "/home")
	})

	r.Run(":8080")
}
