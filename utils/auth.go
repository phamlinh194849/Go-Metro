package utils

import (
  "crypto/md5"
  "fmt"
  "time"

  "github.com/gin-gonic/gin"
  "github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

// Claims represents JWT claims
type Claims struct {
  UserID   uint   `json:"user_id"`
  Username string `json:"username"`
  Role     int    `json:"role"`
  jwt.RegisteredClaims
}

// HashPassword creates MD5 hash of password (đơn giản cho demo)
func HashPassword(password string) string {
  hash := md5.Sum([]byte(password))
  return fmt.Sprintf("%x", hash)
}

// GenerateToken creates JWT token
func GenerateToken(userID uint, username string, role int) (string, error) {
  claims := Claims{
    UserID:   userID,
    Username: username,
    Role:     role,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
      IssuedAt:  jwt.NewNumericDate(time.Now()),
    },
  }

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  return token.SignedString(jwtSecret)
}

// ValidateToken validates JWT token
func ValidateToken(tokenString string) (*Claims, error) {
  token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
    return jwtSecret, nil
  })

  if err != nil {
    return nil, err
  }

  if claims, ok := token.Claims.(*Claims); ok && token.Valid {
    return claims, nil
  }

  return nil, fmt.Errorf("invalid token")
}

// AuthMiddleware middleware để xác thực JWT
func AuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    tokenString := c.GetHeader("Authorization")
    if tokenString == "" {
      c.JSON(401, gin.H{"error": "Authorization header required"})
      c.Abort()
      return
    }

    // Remove "Bearer " prefix if present
    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
      tokenString = tokenString[7:]
    }

    claims, err := ValidateToken(tokenString)
    if err != nil {
      c.JSON(401, gin.H{"error": "Invalid token"})
      c.Abort()
      return
    }

    // Set user info in context
    c.Set("user_id", claims.UserID)
    c.Set("username", claims.Username)
    c.Set("role", claims.Role)

    c.Next()
  }
}

// AdminMiddleware middleware để kiểm tra role admin
func AdminMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists {
      c.JSON(401, gin.H{"error": "User not authenticated"})
      c.Abort()
      return
    }

    if role != 1 {
      c.JSON(403, gin.H{"error": "Admin access required"})
      c.Abort()
      return
    }

    c.Next()
  }
}
