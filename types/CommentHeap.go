package types

import "container/heap"

type CommentHeap []Comment

func (h *CommentHeap) Len() int {
	return len(*h)
}

func (h CommentHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h CommentHeap) Less(i, j int) bool {
	return h[i].Likes > h[j].Likes
}

func (h *CommentHeap) Push(x interface{}) {
	*h = append(*h, x.(Comment))
}

func (h *CommentHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Gets the top 100 comments by likes via PriortyQueue
func GetTopCommentsByLikes(comments []Comment) []Comment {

	TopComments := make([]Comment, 0)
	pq := make(CommentHeap, 0)
	heap.Init(&pq)

	for _, comment := range comments {
		heap.Push(&pq, comment)
	}

	count := 0
	for pq.Len() > 0 && count < 100 {
		TopComments = append(TopComments, heap.Pop(&pq).(Comment))
		count++
	}

	return TopComments
}
