package pagination

import (
    "strconv"
    "github.com/gin-gonic/gin"
)

type Params struct {
    Page  int
    Limit int
}

func FromContext(c *gin.Context) Params {
    page,  _ := strconv.Atoi(c.DefaultQuery("page",  "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    if page  < 1   { page  = 1 }
    if limit < 1   { limit = 10 }
    if limit > 100 { limit = 100 }

    return Params{Page: page, Limit: limit}
}

func (p Params) Offset() int {
    return (p.Page - 1) * p.Limit
}