package main

import (
	"fmt"
	NotionClient "github.com/Jasondeepny/weread-notionso/notion"
	WeRead "github.com/Jasondeepny/weread-notionso/wxread"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Environment variables not fetched : %s", err)
	}

	WeRead.WxCookie = os.Getenv("WEREAD_TOKEN")
	NotionClient.NotionToken = os.Getenv("NOTION_TOKEN")
	NotionClient.NotionDataBaseId = os.Getenv("NOTION_DATABASE_ID")
	log.Printf("WEREAD_TOKEN IS >>>: %s,"+
		"\nNOTION_TOKEN IS >>>: %s,"+
		"\nNOTION_DATABASE_ID IS >>>: %s",
		WeRead.WxCookie, NotionClient.NotionToken, NotionClient.NotionDataBaseId)

	//request WeRead home page
	WeRead.LoadHomePage()

	//Get latest sort number
	latestSort := NotionClient.GetLatestSort()
	fmt.Println(latestSort)

	//Get notebook list
	notebooks := NotionClient.GetNotebookList(WeRead.WereadNotebooksUrl, "GET")
	//fmt.Println(notebooks)

	//Process each notebook
	for _, book := range notebooks.Books {
		bookSort := book.Sort
		//if bookSort <= latestSort {
		//	continue
		//}

		//get bookInfo
		bookInfo := NotionClient.GetBookInfo(WeRead.WereadBookInfo, "GET", book.BookId)

		//check book
		NotionClient.CheckBook(book.BookId)

		//get Chapter data
		chapter := NotionClient.GetChapterInfo(book.BookId)

		//get Bookmark data
		bookmarkList := NotionClient.GetBookmarkList(WeRead.WereadBookmarklistUrl, "GET", book.BookId)

		//get Review data
		summary, reviewList := NotionClient.GetReviewList(WeRead.WereadReviewListUrl, "GET", book.BookId)
		for _, review := range *reviewList {
			bookmarkList.BookMarkUpdates = append(bookmarkList.BookMarkUpdates, review.Review)
		}

		//sort bookmarkList
		NotionClient.SortBooks(bookmarkList.BookMarkUpdates)

		//get Children block
		children, grandChild := NotionClient.GetChildren(chapter, summary, bookmarkList)

		//Insert data into Notion
		id := NotionClient.InsertToNotion(book.Book.Title, book.BookId, book.Book.Cover, bookSort, book.Book.Author, bookInfo.Isbn, bookInfo.NewRating)

		//add Children block
		results := NotionClient.AddChildren(id, children)
		if len(*grandChild) > 0 && results != nil {
			NotionClient.AddGrandchild(grandChild, results)
		}
	}
	fmt.Printf("Synchronization task exits after execution ......")
}
