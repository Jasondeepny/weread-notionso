package notion

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Jasondeepny/notionapi"
	WeRead "github.com/Jasondeepny/weread-notionso/wxread"
	"log"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	NotionDataBaseId string
	NotionToken      string
)

type Book struct {
	BookId         string     `json:"bookId"`
	Title          string     `json:"title"`
	Author         string     `json:"author"`
	Cover          string     `json:"cover"`
	Version        int        `json:"version"`
	Format         string     `json:"format"`
	Type           int        `json:"type"`
	Price          float64    `json:"price"`
	OriginalPrice  int        `json:"originalPrice"`
	Soldout        int        `json:"soldout"`
	BookStatus     int        `json:"bookStatus"`
	PayType        int        `json:"payType"`
	CentPrice      int        `json:"centPrice"`
	Finished       int        `json:"finished"`
	MaxFreeChapter int        `json:"maxFreeChapter"`
	Free           int        `json:"free"`
	MCardDiscount  int        `json:"mcardDiscount"`
	Ispub          int        `json:"ispub"`
	ExtraType      int        `json:"extra_type"`
	Cpid           int        `json:"cpid"`
	PublishTime    string     `json:"publishTime"`
	Categories     []Category `json:"categories"`
	AuthorVids     string     `json:"authorVids"`
	HasLecture     int        `json:"hasLecture"`
	LastChapterIdx int        `json:"lastChapterIdx"`
	BlockSaveImg   int        `json:"blockSaveImg"`
	Language       string     `json:"language"`
	HideUpdateTime bool       `json:"hideUpdateTime"`
}

type BookInfo struct {
	BookId             string  `json:"bookId"`
	Book               Book    `json:"book"`
	ReviewCount        int     `json:"reviewCount"`
	ReviewLikeCount    int     `json:"reviewLikeCount"`
	ReviewCommentCount int     `json:"reviewCommentCount"`
	NoteCount          int     `json:"noteCount"`
	BookmarkCount      int     `json:"bookmarkCount"`
	Sort               float64 `json:"sort"`
}

type BookInfoData struct {
	BookId    string  `json:"bookId"`
	Isbn      string  `json:"isbn"`
	NewRating float64 `json:"newRating"`
}

type BookList struct {
	SyncKey           int        `json:"synckey"`
	TotalBookCount    int        `json:"totalBookCount"`
	NoBookReviewCount int        `json:"noBookReviewCount"`
	Books             []BookInfo `json:"books"`
}

type Category struct {
	CategoryId    int    `json:"categoryId"`
	SubCategoryId int    `json:"subCategoryId"`
	CategoryType  int    `json:"categoryType"`
	Title         string `json:"title"`
}

type ChapterData struct {
	Data []BookData `json:"data"`
}

type ReviewData struct {
	Reviews []Reviews `json:"reviews"`
}
type Reviews struct {
	Review BookMarkUpdate `json:"review"`
}

type BookData struct {
	BookId  string    `json:"bookId"`
	Updates []Updated `json:"updated"`
}

type Updated struct {
	ChapterUid int    `json:"chapterUid"`
	ChapterIdx int    `json:"chapterIdx"`
	UpdateTime int    `json:"updateTime"`
	ReadAhead  int    `json:"readAhead"`
	Title      string `json:"title"`
	Level      int    `json:"level"`
}

type BookMarkUpdate struct {
	BookId     string `json:"bookId"`
	Range      string `json:"range"`
	MarkText   string `json:"markText"`
	Content    string `json:"content"`
	ReviewId   string `json:"reviewId"`
	Abstract   string `json:"abstract"`
	Type       int    `json:"type"`
	ChapterUid int    `json:"chapterUid"`
	CreateTime int64  `json:"createTime"`
	Style      int    `json:"style"`
	ColorStyle int    `json:"colorStyle"`
}

type ReadInfo struct {
	FinishedBookIndex    int    `json:"finishedBookIndex"`
	FinishedDate         int64  `json:"finishedDate"`
	ReadingBookCount     int    `json:"readingBookCount"`
	ReadingBookDate      int64  `json:"readingBookDate"`
	ReadingProgress      int    `json:"readingProgress"`
	ReadingReviewId      string `json:"readingReviewId"`
	MarkedStatus         int    `json:"markedStatus"`
	ReadingTime          int64  `json:"readingTime"`
	TotalReadDay         int    `json:"totalReadDay"`
	RecordReadingTime    int64  `json:"recordReadingTime"`
	DeepestNightReadTime int64  `json:"deepestNightReadTime"`
	ContinueReadDays     int    `json:"continueReadDays"`
	ContinueBeginDate    int64  `json:"continueBeginDate"`
	ContinueEndDate      int64  `json:"continueEndDate"`
	ShowSummary          int    `json:"showSummary"`
	ShowDetail           int    `json:"showDetail"`
}

type BookMarks struct {
	BookMarkUpdates []BookMarkUpdate `json:"updated"`
	//Chapters        []struct {
	//	BookId     string `json:"bookId"`
	//	ChapterUid int    `json:"chapterUid"`
	//	ChapterIdx int    `json:"chapterIdx"`
	//	Title      string `json:"title"`
	//} `json:"chapters"`
}

func Client() *notionapi.Client {
	return notionapi.NewClient(notionapi.Token(NotionToken))
}

func getHeading(children *notionapi.Blocks, level int, content string) {
	if level == 1 {
		*children = append(*children, &notionapi.Heading1Block{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeHeading1,
			},
			Heading1: struct {
				RichText     []notionapi.RichText `json:"rich_text"`
				Children     notionapi.Blocks     `json:"children,omitempty"`
				Color        string               `json:"color,omitempty"`
				IsToggleable bool                 `json:"is_toggleable,omitempty"`
			}{[]notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{Content: content},
				},
			}, nil, "default", false,
			},
		})
	} else if level == 2 {
		*children = append(*children, &notionapi.Heading2Block{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeHeading2,
			},
			Heading2: struct {
				RichText     []notionapi.RichText `json:"rich_text"`
				Children     notionapi.Blocks     `json:"children,omitempty"`
				Color        string               `json:"color,omitempty"`
				IsToggleable bool                 `json:"is_toggleable,omitempty"`
			}{[]notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{Content: content},
				},
			}, nil, "default", false,
			},
		})
	} else {
		*children = append(*children, &notionapi.Heading3Block{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeHeading3,
			},
			Heading3: struct {
				RichText     []notionapi.RichText `json:"rich_text"`
				Children     notionapi.Blocks     `json:"children,omitempty"`
				Color        string               `json:"color,omitempty"`
				IsToggleable bool                 `json:"is_toggleable,omitempty"`
			}{[]notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{Content: content},
				},
			}, nil, "default", false,
			},
		})
	}
}

func getQuote(lens int, grandChildren *map[int]notionapi.Block, content string) {
	n := &notionapi.QuoteBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockQuote,
		},
		Quote: struct {
			RichText []notionapi.RichText `json:"rich_text"`
			Children notionapi.Blocks     `json:"children,omitempty"`
			Color    string               `json:"color,omitempty"`
		}{[]notionapi.RichText{
			{
				Type: notionapi.ObjectTypeText,
				Text: &notionapi.Text{Content: content},
			},
		}, nil, "default",
		},
	}
	(*grandChildren)[lens] = n
}

func getCallout(children *notionapi.Blocks, content string, style int, colorStyle int, reviewId string) {
	//根据不同的划线样式设置不同的emoji 直线type=0 背景颜色是1 波浪线是2
	var emoji notionapi.Emoji = "🌟"
	if style == 0 {
		emoji = "💡"
	} else if style == 1 {
		emoji = "⭐"
	}
	//如果reviewId不是空说明是笔记
	if reviewId != "" {
		emoji = "✍️"
	}
	//根据划线颜色设置文字的颜色
	var color notionapi.Color = "default"
	if colorStyle == 1 {
		color = notionapi.ColorRed
	} else if colorStyle == 2 {
		color = notionapi.ColorPurple
	} else if colorStyle == 3 {
		color = notionapi.ColorBlue
	} else if colorStyle == 4 {
		color = notionapi.ColorGreen
	} else if colorStyle == 5 {
		color = notionapi.ColorYellow
	}
	*children = append(*children, &notionapi.CalloutBlock{
		BasicBlock: notionapi.BasicBlock{
			Object: notionapi.ObjectTypeBlock,
			Type:   notionapi.BlockCallout,
		},
		Callout: notionapi.Callout{
			RichText: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: &notionapi.Text{
						Content: content,
					},
				},
			},
			Icon: &notionapi.Icon{
				Type:  "emoji",
				Emoji: &emoji,
			},
			Color: string(color),
		},
	})
}

func InsertToNotion(bookName string, bookId string, cover string, sort float64, author string, isbn string, rating float64) string {
	//"""插入到notion"""
	parent := notionapi.Parent{
		Type:       notionapi.ParentTypeDatabaseID,
		DatabaseID: notionapi.DatabaseID(NotionDataBaseId),
	}
	properties := notionapi.Properties{
		"BookName": notionapi.TitleProperty{Title: []notionapi.RichText{{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: bookName}}}},
		"BookId":   notionapi.RichTextProperty{RichText: []notionapi.RichText{{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: bookId}}}},
		"ISBN":     notionapi.RichTextProperty{RichText: []notionapi.RichText{{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: isbn}}}},
		"URL":      notionapi.URLProperty{URL: fmt.Sprintf("https://weread.qq.com/web/reader/%s", calculateBookStrId(bookId))},
		"Author":   notionapi.RichTextProperty{RichText: []notionapi.RichText{{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: author}}}},
		"Sort":     notionapi.NumberProperty{Number: sort},
		"Rating":   notionapi.NumberProperty{Number: rating},
		"Cover":    notionapi.FilesProperty{Files: []notionapi.File{{Type: notionapi.FileTypeExternal, Name: "Cover", External: &notionapi.FileObject{URL: cover}}}},
	}
	readInfo := getReadInfo(bookId)
	if readInfo != nil {
		markedStatus := readInfo.MarkedStatus
		readingTime := readInfo.ReadingTime
		formatTime := ""
		hour := readingTime / 3600
		if hour > 0 {
			formatTime += fmt.Sprintf("%d时", hour)
		}
		minutes := (readingTime % 3600) / 60
		if minutes > 0 {
			formatTime += fmt.Sprintf("%d分", minutes)
		}
		properties["Status"] = notionapi.SelectProperty{Select: notionapi.Option{Name: getStatus(markedStatus)}}
		properties["ReadingTime"] = notionapi.RichTextProperty{RichText: []notionapi.RichText{{Type: notionapi.ObjectTypeText, Text: &notionapi.Text{Content: formatTime}}}}
		if readInfo.FinishedDate != 0 {
			properties["Date"] = notionapi.DateProperty{Date: &notionapi.DateObject{
				Start: (*notionapi.Date)(timeFormat(readInfo.FinishedDate)),
			}}
		}
	}
	icon := &notionapi.Icon{Type: notionapi.FileTypeExternal, External: &notionapi.FileObject{URL: cover}}
	page, err := Client().Page.Create(context.Background(), &notionapi.PageCreateRequest{
		Parent:     parent,
		Properties: properties,
		Icon:       icon,
	})
	check(err)
	return string(page.ID)
}

func timeFormat(date int64) *time.Time {
	t := time.Unix(date, 0)
	return &t
}

func getStatus(markedStatus int) string {
	if markedStatus == 4 {
		return "读完"
	} else {
		return "在读"
	}
}

func AddChildren(id string, children *notionapi.Blocks) *notionapi.Blocks {
	result := make(notionapi.Blocks, 0)
	for i := 0; i < len(*children)/100+1; i++ {
		time.Sleep(300 * time.Millisecond)
		startIndex, endIndex := i*100, (i+1)*100
		if endIndex > len(*children) {
			endIndex = len(*children)
		}
		res, err := Client().Block.AppendChildren(context.Background(), notionapi.BlockID(id), &notionapi.AppendBlockChildrenRequest{
			//After:    notionapi.BlockID(id),
			Children: (*children)[startIndex:endIndex],
		})
		check(err)
		if res.Results != nil && len(res.Results) > 0 {
			result = append(result, res.Results...)
		}
	}
	return &result
}

func AddGrandchild(grandchild *map[int]notionapi.Block, results *notionapi.Blocks) {
	for k, v := range *grandchild {
		//blockID := (*results)[k].GetID()
		_, err := Client().Block.AppendChildren(context.Background(), (*results)[k].GetID(), &notionapi.AppendBlockChildrenRequest{
			//After:    blockID,
			Children: []notionapi.Block{v},
		})
		check(err)
	}
}

func GetChildren(chapter *[]BookData, reviews *[]Reviews, bookMarks *BookMarks) (*notionapi.Blocks, *map[int]notionapi.Block) {
	d := make(map[string][]BookMarkUpdate)
	children := make(notionapi.Blocks, len(*chapter))
	grandchild := make(map[int]notionapi.Block, 8)
	bookDataMap := make(map[string][]Updated, 8)
	if chapter != nil && len(*chapter) > 0 {
		children[0] = &notionapi.TableOfContentsBlock{
			BasicBlock: notionapi.BasicBlock{
				Object: notionapi.ObjectTypeBlock,
				Type:   notionapi.BlockTypeTableOfContents,
			},
			TableOfContents: notionapi.TableOfContents{Color: "default"},
		}
		for _, bookmark := range bookMarks.BookMarkUpdates {
			uid := strconv.Itoa(bookmark.ChapterUid)
			d[uid] = append(d[uid], bookmark)
		}
		for _, ch := range *chapter {
			for _, update := range ch.Updates {
				uid := strconv.Itoa(update.ChapterUid)
				if _, exist := bookDataMap[uid]; !exist {
					bookDataMap[uid] = append(bookDataMap[uid], update)
				}
			}
		}
		for k, v := range d {
			if _, ok := bookDataMap[k]; ok {
				getHeading(&children, bookDataMap[k][0].Level, bookDataMap[k][0].Title)
			}
			for _, s := range v {
				getCallout(&children, s.MarkText, s.Style, s.ColorStyle, s.ReviewId)
				if s.Abstract != "" {
					//fmt.Println(fmt.Sprintf("children‘s lens is %d", len(children)-1))
					getQuote(len(children)-1, &grandchild, s.Abstract)
				}
			}
		}
	} else {
		for _, s := range bookMarks.BookMarkUpdates {
			getCallout(&children, s.MarkText, s.Style, s.ColorStyle, s.ReviewId)
		}
	}

	if reviews != nil && len(*reviews) > 0 {
		getHeading(&children, 1, "点评")
		for _, r := range *reviews {
			getCallout(&children, r.Review.Content, r.Review.Style, r.Review.ColorStyle, r.Review.ReviewId)
		}
	}
	return &children, &grandchild
}

func GetBookInfo(Url string, method string, bookId string) *BookInfoData {
	params := url.Values{}
	params.Add("bookId", bookId)
	Url = fmt.Sprintf("%s?%s", Url, params.Encode())
	res := WeRead.DoWxQuery(Url, method, nil)
	var bookInfo BookInfoData
	err := json.Unmarshal(res, &bookInfo)
	check(err)
	return &bookInfo
}

func SortBooks(bookmarkList []BookMarkUpdate) {
	sort.Slice(bookmarkList, func(i, j int) bool {
		// Sort by chapterUid
		if bookmarkList[i].ChapterUid != bookmarkList[j].ChapterUid {
			return bookmarkList[i].ChapterUid < bookmarkList[j].ChapterUid
		}
		// Sort by range
		rangeI := bookmarkList[i].Range
		rangeJ := bookmarkList[j].Range
		// Handle empty range
		if rangeI == "" && rangeJ == "" {
			return false
		} else if rangeI == "" {
			return true
		} else if rangeJ == "" {
			return false
		}
		// Split range string and convert start value to integer
		rangeISplit := strings.Split(rangeI, "-")
		rangeJSplit := strings.Split(rangeJ, "-")
		rangeIStart, _ := strconv.Atoi(rangeISplit[0])
		rangeJStart, _ := strconv.Atoi(rangeJSplit[0])
		return rangeIStart < rangeJStart
	})
}

func GetReviewList(Url string, method string, bookId string) (*[]Reviews, *[]Reviews) {
	params := url.Values{}
	params.Add("bookId", bookId)
	params.Add("listType", "11")
	params.Add("mine", "1")
	params.Add("syncKey", "0")
	Url = fmt.Sprintf("%s?%s", Url, params.Encode())
	res := WeRead.DoWxQuery(Url, method, nil)
	var reviewData ReviewData
	err := json.Unmarshal(res, &reviewData)
	check(err)
	var summary []Reviews
	var reviews []Reviews
	for _, review := range reviewData.Reviews {
		if review.Review.Type == 4 {
			summary = append(summary, review)
		}
		if review.Review.Type == 1 {
			review.Review.MarkText = review.Review.Content
			reviews = append(reviews, review)
		}
	}
	return &summary, &reviews
}

func getReadInfo(bookId string) *ReadInfo {
	params := url.Values{}
	params.Add("bookId", bookId)
	params.Add("readingDetail", "1")
	params.Add("readingBookIndex", "1")
	params.Add("finishedDate", "1")
	Url := fmt.Sprintf("%s?%s", WeRead.WereadReadInfoUrl, params.Encode())
	res := WeRead.DoWxQuery(Url, "GET", nil)
	var readInfo ReadInfo
	err := json.Unmarshal(res, &readInfo)
	check(err)
	return &readInfo
}

func calculateBookStrId(bookId string) string {
	md5Hash := md5.New()
	md5Hash.Write([]byte(bookId))
	digest := hex.EncodeToString(md5Hash.Sum(nil))
	result := digest[:3]

	code, transformedIds := transformId(bookId)
	result += code + "2" + digest[len(digest)-2:]

	for i := 0; i < len(transformedIds); i++ {
		hexLengthStr := fmt.Sprintf("%x", len(transformedIds[i]))
		if len(hexLengthStr) == 1 {
			hexLengthStr = "0" + hexLengthStr
		}

		result += hexLengthStr + transformedIds[i]

		if i < len(transformedIds)-1 {
			result += "g"
		}
	}

	if len(result) < 20 {
		result += digest[:20-len(result)]
	}

	md5Hash = md5.New()
	md5Hash.Write([]byte(result))
	result += hex.EncodeToString(md5Hash.Sum(nil))[:3]
	return result
}

func transformId(bookId string) (string, []string) {
	idLength := len(bookId)
	if matched, _ := regexp.MatchString(`^\d*$`, bookId); matched {
		var transformedIds []string
		for i := 0; i < idLength; i += 9 {
			end := i + 9
			if end > idLength {
				end = idLength
			}
			ary := fmt.Sprintf("%x", toInt(bookId[i:end]))
			transformedIds = append(transformedIds, ary)
		}
		return "3", transformedIds
	}

	var result strings.Builder
	for i := 0; i < idLength; i++ {
		result.WriteString(fmt.Sprintf("%x", int(bookId[i])))
	}
	return "4", []string{result.String()}
}

func toInt(s string) int {
	var result int
	for _, r := range s {
		result = result*10 + int(r-'0')
	}
	return result
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetBookmarkList(Url string, method string, bookId string) *BookMarks {
	params := url.Values{}
	params.Add("bookId", bookId)
	Url = fmt.Sprintf("%s?%s", Url, params.Encode())
	res := WeRead.DoWxQuery(Url, method, nil)
	var bookMarks BookMarks
	err := json.Unmarshal(res, &bookMarks)
	check(err)
	SortBooks(bookMarks.BookMarkUpdates)
	return &bookMarks
}

func GetChapterInfo(bookId string) *[]BookData {
	requestBody := map[string]interface{}{
		"bookIds":  []string{bookId},
		"synckeys": []int{0},
		"teenmode": 0,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	check(err)
	res := WeRead.DoWxQuery(WeRead.WereadChapterInfo, "POST", requestBodyBytes)
	var chapterData ChapterData
	err = json.Unmarshal(res, &chapterData)
	check(err)
	return &chapterData.Data
}

func CheckBook(bookId string) {
	time.Sleep(300 * time.Millisecond)
	response, err := Client().Database.Query(context.Background(), notionapi.DatabaseID(NotionDataBaseId), &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "BookId",
			RichText: &notionapi.TextFilterCondition{
				Equals: bookId,
			},
		},
	})
	check(err)
	if nil != response.Results {
		for _, r := range response.Results {
			time.Sleep(1)
			_, err := Client().Block.Delete(context.Background(), notionapi.BlockID(r.ID))
			check(err)
		}
	}
}

func GetNotebookList(Url string, method string) BookList {
	var books BookList
	res := WeRead.DoWxQuery(Url, method, nil)
	err := json.Unmarshal(res, &books)
	check(err)
	return books
}

func GetLatestSort() float64 {
	response, err := Client().Database.Query(context.Background(), notionapi.DatabaseID(NotionDataBaseId), &notionapi.DatabaseQueryRequest{
		Filter: &notionapi.PropertyFilter{
			Property: "Sort",
			Number: &notionapi.NumberFilterCondition{
				IsNotEmpty: true,
			},
		},
		Sorts: []notionapi.SortObject{
			{
				Property:  "Sort",
				Direction: "descending",
			},
		},
		PageSize: 1,
	})
	if err != nil {
		fmt.Printf("Failed to query database: %s\n", err.Error())
		return 0
	}
	fmt.Printf("Query results: %+v\n", response.Results)
	if len(response.Results) == 1 {
		property := response.Results[0].Properties["Sort"]
		if numberProp, ok := property.(*notionapi.NumberProperty); ok {
			sortValue := numberProp.Number
			fmt.Printf("Sort value: %v\n", sortValue)
			return sortValue
		}
	}
	return 0
}
