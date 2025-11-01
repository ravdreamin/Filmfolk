package handlers

import (
	"net/http"
	"strconv"

	"filmfolk/internal/middleware"
	"filmfolk/internal/services"

	"github.com/gin-gonic/gin"
)

type FollowerHandler struct {
	followerService *services.FollowerService
}

func NewFollowerHandler() *FollowerHandler {
	return &FollowerHandler{
		followerService: services.NewFollowerService(),
	}
}

// FollowUser handles POST /api/v1/users/:id/follow
func (h *FollowerHandler) FollowUser(c *gin.Context) {
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followerID := middleware.GetUserID(c)

	if err := h.followerService.FollowUser(followerID, followingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully followed user"})
}

// UnfollowUser handles DELETE /api/v1/users/:id/follow
func (h *FollowerHandler) UnfollowUser(c *gin.Context) {
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followerID := middleware.GetUserID(c)

	if err := h.followerService.UnfollowUser(followerID, followingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully unfollowed user"})
}

// CheckFollowStatus handles GET /api/v1/users/:id/follow/status
func (h *FollowerHandler) CheckFollowStatus(c *gin.Context) {
	followingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	followerID := middleware.GetUserID(c)

	isFollowing, err := h.followerService.IsFollowing(followerID, followingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_following": isFollowing,
	})
}

// GetFollowers handles GET /api/v1/users/:id/followers
func (h *FollowerHandler) GetFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	followers, total, err := h.followerService.GetFollowers(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetFollowing handles GET /api/v1/users/:id/following
func (h *FollowerHandler) GetFollowing(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	following, total, err := h.followerService.GetFollowing(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"following": following,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
