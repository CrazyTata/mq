package assembler

import (
	domain "mq/domain/qa"
	external "mq/infrastructure/integration/qa"
)

func EOTODOChatSearch(chatSearch *external.ChatSearch) *domain.ChatSearch {
	return &domain.ChatSearch{
		Message:      chatSearch.Message,
		HistoryID:    chatSearch.HistoryID,
		SessionID:    chatSearch.SessionID,
		Number:       chatSearch.Number,
		IsEnd:        chatSearch.IsEnd,
		IsSearch:     chatSearch.IsSearch,
		ErrorMessage: chatSearch.ErrorMessage,
		Code:         chatSearch.Code,
	}
}

func EOTODOGetHistory(history []*external.History) []*domain.History {
	historyList := make([]*domain.History, 0, len(history))
	for _, h := range history {
		historyList = append(historyList, &domain.History{
			ChatTopic: h.ChatTopic,
			SessionID: h.SessionID,
			UpdateAt:  h.UpdateAt,
		})
	}
	return historyList
}

func EOTODOGetHistoryDetail(historyDetail []*external.HistoryDetail) []*domain.HistoryDetail {
	historyDetailList := make([]*domain.HistoryDetail, 0, len(historyDetail))
	for _, h := range historyDetail {
		historyDetailList = append(historyDetailList, &domain.HistoryDetail{
			HistoryID: h.HistoryID,
			ID:        h.ID,
			Message:   h.Message,
			Self:      h.Self,
		})
	}
	return historyDetailList
}

func EOTODOGetRetrievalDocs(retrievalDocs []*external.RetrievalDoc) []*domain.RetrievalDoc {
	retrievalDocsList := make([]*domain.RetrievalDoc, 0, len(retrievalDocs))
	for _, r := range retrievalDocs {
		retrievalDocsList = append(retrievalDocsList, &domain.RetrievalDoc{
			ChuckID:      r.ChuckID,
			ShortContent: r.ShortContent,
			SourceDoc:    r.SourceDoc,
			Topic:        r.Topic,
			Type:         r.Type,
			URL:          r.URL,
		})
	}
	return retrievalDocsList
}
