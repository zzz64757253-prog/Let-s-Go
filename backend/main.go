package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sashabaranov/go-openai"
)

var dbPool *pgxpool.Pool

func initDB() {
	var err error
	dbPool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	err = dbPool.Ping(context.Background())
	if err != nil {

		log.Fatal("เชื่อมต่อ Database ไม่สำเร็จ:", err)
	}
	log.Println("เชื่อมต่อ Database สำเร็จจริงแล้ว!")
}

var (
	r2           *minio.Client
	bucketName   string
	publicDomain string
)

func initFile() {
	var err error
	accountID := os.Getenv("R2_ENDPOINT")
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")

	r2, err = minio.New(accountID+".r2.cloudflarestorage.com", &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		panic("เชื่อมต่อ R2 ไม่สำเร็จ: " + err.Error())
	}
}

type IncomingMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatRequest struct {
	Messages []openai.ChatCompletionMessage `json:"messages"`
	Model    string
}
type MenuType struct {
	Menu_id   int     `json:"menu_id"`
	Name_menu string  `json:"name_menu"`
	Menu_url  *string `json:"menu_url"`
	Parent_id *int    `json:"parent_id"`
	Menu_icon *string `json:"menu_icon"`
}
type PostType struct {
	Post_uid  string     `json:"post_uid"`
	Username  string     `json:"username"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Image_url *string    `json:"image_url"`
	Post_date *time.Time `json:"post_date"`
	Pin       string     `json:"pin"`
}
type ModelType struct {
	Modelname string `json:"modelname"`
}
type SkillType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Score       int    `json:"score"`
	Image_url   string `json:"image_url"`
	Skill_type  int    `json:"skill_type"`
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("โหลด .env ไม่สำเร็จ:", err)
	}

	initDB()
	defer dbPool.Close()

	initFile()

	r := gin.Default()

	r.GET("/api/lottery", func(c *gin.Context) {

		payload := []byte(`{}`)
		req, err := http.NewRequest("POST", "https://www.glo.or.th/api/lottery/getLatestLottery", bytes.NewBuffer(payload))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(resp.StatusCode, "application/json", bodyBytes)
	})

	config := openai.DefaultConfig(os.Getenv("GEMINI_API_KEY"))
	config.BaseURL = "https://generativelanguage.googleapis.com/v1beta/openai"
	client := openai.NewClientWithConfig(config)
	r.POST("/api/chat", func(c *gin.Context) {
		var req ChatRequest
		c.ShouldBindJSON(&req)
		messages := append([]openai.ChatCompletionMessage{
			{Role: "system", Content: "`คุณคือผู้ช่วยเว็บ Portfolio ของ อัครพล (Zee) นักพัฒนา Full-stack (Nuxt, Go) หน้าที่: แนะนำทักษะ ผลงาน และให้ข้อมูลจ้างงาน สนใจจ้างงานติดต่อ: - Line: [ไอดีไลน์] - Tel: [เบอร์โทร] - Email: [อีเมล] ตอบสุภาพ สั้นกระชับ ตรงประเด็น`"},
		}, req.Messages...)
		resp, err := client.CreateChatCompletion(c.Request.Context(), openai.ChatCompletionRequest{
			Model:    req.Model,
			Messages: messages,
		})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})

			return
		}
		c.JSON(200, gin.H{"message": resp.Choices[0].Message.Content})
	})

	r.POST("/api/upload", func(c *gin.Context) {
		bucketName = os.Getenv("R2_BUCKET_NAME")
		publicDomain = os.Getenv("R2_PUBLIC_DOMAIN")

		fileHeader, err := c.FormFile("file")
		if err != nil {
			fmt.Println("111")
			c.JSON(http.StatusBadRequest, gin.H{"statusMessage": "กรุณาเลือกไฟล์ที่ต้องการอัปโหลด"})
			return
		}
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println("222")
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "อ่านไฟล์ไม่สำเร็จ"})
			return
		}
		defer file.Close()
		ext := filepath.Ext(fileHeader.Filename)
		newFilename := fmt.Sprintf("%d-%x%s", time.Now().UnixNano(), rand.Int31(), ext)
		fmt.Println(newFilename)
		fmt.Println(bucketName)
		fmt.Println(fileHeader.Size)
		fmt.Println(fileHeader.Header.Get("Content-Type"))
		_, err = r2.PutObject(c.Request.Context(), bucketName, newFilename, file, fileHeader.Size, minio.PutObjectOptions{
			ContentType: fileHeader.Header.Get("Content-Type"),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"statusMessage": "อัปโหลดขึ้น R2 ไม่สำเร็จ: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"url":      fmt.Sprintf("%s/%s", publicDomain, newFilename),
			"filename": newFilename,
		})
	})

	r.DELETE("/api/delete", func(c *gin.Context) {
		bucketName = os.Getenv("R2_BUCKET_NAME")
		var req struct {
			URL string `json:"url" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": "failure",
				"message": "ลบไฟล์ไม่สำเร็จ: " + err.Error()})
			return
		}
		fileKey := path.Base(req.URL)

		err := r2.RemoveObject(c.Request.Context(), bucketName, fileKey, minio.RemoveObjectOptions{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": "failure",
				"message": "ลบไฟล์ไม่สำเร็จ: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": "success",
			"message": "ลบภาพเรียบร้อยแล้ว",
		})
	})

	r.GET("/api/menu", func(c *gin.Context) {
		var mainMenus []MenuType
		var subMenus []MenuType
		rows, err := dbPool.Query(context.Background(), "SELECT menu_id, name_menu, menu_url, parent_id, menu_icon FROM menu WHERE parent_id IS NULL ORDER BY menu_id ASC")
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var m MenuType
			if err := rows.Scan(&m.Menu_id, &m.Name_menu, &m.Menu_url, &m.Parent_id, &m.Menu_icon); err != nil {
				return
			}
			mainMenus = append(mainMenus, m)
		}

		rows, err = dbPool.Query(context.Background(), "SELECT menu_id, name_menu, menu_url, parent_id, menu_icon FROM menu WHERE parent_id IS NOT NULL ORDER BY menu_id ASC")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rows.Close()

		for rows.Next() {
			var sm MenuType
			if err := rows.Scan(&sm.Menu_id, &sm.Name_menu, &sm.Menu_url, &sm.Parent_id, &sm.Menu_icon); err != nil {
				fmt.Println(err.Error())
				return
			}
			subMenus = append(subMenus, sm)
		}
		c.JSON(http.StatusOK, gin.H{
			"mainMenus": mainMenus,
			"subMenus":  subMenus,
		})
	})
	r.GET("/api/content", func(c *gin.Context) {
		var posts []PostType
		rows, err := dbPool.Query(context.Background(), "SELECT post_uid, username, title, content, image_url, post_date FROM posts ORDER BY post_date DESC")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rows.Close()
		for rows.Next() {
			var p PostType
			if err := rows.Scan(&p.Post_uid, &p.Username, &p.Title, &p.Content, &p.Image_url, &p.Post_date); err != nil {
				fmt.Println(err.Error())
				return
			}
			posts = append(posts, p)
		}
		//fmt.Println(posts)
		c.JSON(http.StatusOK, gin.H{
			"posts": posts,
		})
	})
	r.PUT("/api/content", func(c *gin.Context) {
		var req PostType
		c.ShouldBindJSON(&req)
		var p struct {
			Pin string `json:"pin" binding:"required"`
		}
		err := dbPool.QueryRow(context.Background(), "SELECT pin FROM posts WHERE post_uid = $1", req.Post_uid).Scan(&p.Pin)
		if err != nil {
			c.JSON(http.StatusFound, gin.H{
				"success": "failure",
				"message": "err",
			})
			return
		}
		if p.Pin != req.Pin {
			c.JSON(http.StatusFound, gin.H{
				"success": "failure",
				"message": "PIN ไม่ถูกต้อง",
			})
			return
		}
		cmdTag, err := dbPool.Exec(context.Background(), "UPDATE posts SET username = $1, title = $2, content = $3, image_url = $4 WHERE post_uid = $5",
			req.Username,
			req.Title, req.Content,
			req.Image_url,
			req.Post_uid,
		)
		if err != nil {
			return
		}
		if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": "failure",
				"message": "ไม่พบโพสต์ที่ต้องการแก้ไข",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": "success",
			"message": "แก้ไขโพสต์สำเร็จ",
		})
	})
	r.DELETE("/api/content", func(c *gin.Context) {
		var req struct {
			Post_uid string `json:"post_uid" binding:"required"`
			Pin      string `json:"pin" binding:"required"`
		}
		c.ShouldBindJSON(&req)

		cmdTag, err := dbPool.Exec(context.Background(), "DELETE FROM posts WHERE post_uid = $1 AND pin = $2",
			req.Post_uid,
			req.Pin,
		)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": "failure",
				"message": "query err",
			})
			return
		}
		if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": "failure",
				"message": "รหัสไม่ถูกต้อง",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": "success",
			"message": "ลบสำเร็จ!!",
		})
	})
	r.POST("/api/content", func(c *gin.Context) {
		var req PostType
		c.ShouldBindJSON(&req)
		_, err := dbPool.Exec(context.Background(), "INSERT INTO posts (username, title, content, image_url, pin) VALUES ($1, $2, $3, $4, $5)",
			req.Username,
			req.Title,
			req.Content,
			req.Image_url,
			req.Pin,
		)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"success": "failure",
				"message": "เกิดข้อผิดพลาดบางอย่าง!",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": "success",
			"message": "เพิ่มโพสต์สำเร็จ!",
		})

	})
	r.GET("/api/mobilemenu", func(c *gin.Context) {
		var menus []MenuType
		rows, err := dbPool.Query(context.Background(), "SELECT menu_id, name_menu, menu_url, menu_icon FROM menu WHERE menu_url IS NOT NULL ORDER BY menu_id ASC")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var ms MenuType
			if err := rows.Scan(&ms.Menu_id, &ms.Name_menu, &ms.Menu_url, &ms.Menu_icon); err != nil {
				fmt.Println(err)
				return
			}
			menus = append(menus, ms)
		}
		c.JSON(http.StatusOK, gin.H{
			"menus": menus,
		})
	})
	r.GET("/api/model", func(c *gin.Context) {
		var models []ModelType
		rows, err := dbPool.Query(context.Background(), "SELECT modelname FROM models")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var m ModelType
			if err := rows.Scan(&m.Modelname); err != nil {
				fmt.Println(err)
				return
			}
			models = append(models, m)
		}
		c.JSON(http.StatusOK, gin.H{
			"models": models,
		})
	})
	r.GET("/api/skill", func(c *gin.Context) {
		var skills []SkillType
		rows, err := dbPool.Query(context.Background(), "SELECT name, description, score, image_url, skill_type FROM skills ORDER BY  skills.skill_type ASC, skills.order ASC")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var s SkillType
			err := rows.Scan(&s.Name, &s.Description, &s.Score, &s.Image_url, &s.Skill_type)
			if err != nil {
				fmt.Println(err)
				return
			}
			skills = append(skills, s)
		}
		c.JSON(http.StatusOK, skills)
	})
	r.Run(":8080")
}
