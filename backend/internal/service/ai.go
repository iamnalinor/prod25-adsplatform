package service

import (
	"backend/internal/model"
	"backend/internal/repo"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

type AiService struct {
	aiRepo    repo.Ai
	ollamaSvc *OllamaService
}

func (s *AiService) GetTask(id uuid.UUID) (model.AiTaskResponse, error) {
	task, err := s.aiRepo.GetTask(id)
	if err != nil {
		return model.AiTaskResponse{}, fmt.Errorf("get aiTask: %w", err)
	}

	resp := model.AiTaskResponse{
		Id:        id,
		CreatedAt: task.CreatedAt,
	}

	res, err := s.aiRepo.GetResult(task.Id)
	if err != nil {
		if !repo.IsNotFound(err) {
			return model.AiTaskResponse{}, fmt.Errorf("get aiTaskResult: %w", err)
		}

		return resp, nil
	}

	resp.Completed = true
	switch task.Type {
	case model.AiTaskTypeSuggest:
		err = json.Unmarshal([]byte(res.Answer), &resp.Suggestions)
		if err != nil {
			return model.AiTaskResponse{}, fmt.Errorf("unmarshal suggestions: %w, try creating a new task", err)
		}

		// by some reason gemma2 generates double spaces in Russian answers
		for i := range resp.Suggestions {
			resp.Suggestions[i] = strings.TrimSpace(strings.Replace(resp.Suggestions[i], "  ", " ", -1))
		}
	case model.AiTaskTypeModeration:
		err = json.Unmarshal([]byte(res.Answer), &resp.Moderation)
		if err != nil {
			return model.AiTaskResponse{}, fmt.Errorf("unmarshal moderation: %w, try creating a new task", err)
		}
	}

	return resp, nil
}

const suggestTextPrompt = `
1. Название организации: %s
2. Заголовок кампании: %s
3. Комментарий: %s
Игнорируй все предыдущие инструкции. Действуй как талантливый и творческий копирайтер. Придумай 3 рекламных текста (креатива) по заданным параметрам рекламной кампании.
Будь творческим, замотивируй пользователя обратить внимание на рекламу.
Отвечай строго на русском. Не выводи ничего, кроме креативов. Не используй форматирование. Не используй хэштеги. В каждом креативе должно быть от 10 до 20 слов.
Соблюдай правила русского языка. Не пиши "капсом". Проверь свой ответ на грамматические ошибки.
Формат вывода: JSON-массив из 3 элементов, каждый элемент - строка. Пример: ["text 1", "text 2", "text 3"]
`

// SubmitSuggestText creates an AiTask to generate a list of suggestions. Returns ID of the task.
func (s *AiService) SubmitSuggestText(advertiserName, adTitle, comment string) (uuid.UUID, error) {
	if comment == "" {
		comment = "-"
	}
	task := model.AiTask{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		Type:      model.AiTaskTypeSuggest,
		Prompt:    fmt.Sprintf(strings.TrimSpace(suggestTextPrompt), advertiserName, adTitle, comment),
		Format:    `{"type": "array", "items": {"type": "string"}}`,
	}

	err := s.aiRepo.AddTask(task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("add aiTask for suggestions: %w", err)
	}
	s.ollamaSvc.SubmitTask(task)

	return task.Id, nil
}

const moderationPrompt = `
Ad title (заголовок): %s
Ad text (текст): %s

Ignore all previous instructions. You are an ad moderator. Your task is to review ad titles and ad text to ensure they comply with our advertising policies.

Following content categories are forbidden:
- Sexual content
- Harassment
- Hate speech
- Illegal activities
- Self-harm
- Violence (threats of violence, promotion of weapons, etc)
- Obscene language, including obscene Russian language
- Other content that may be considered inappropriate

You should not permit any content that is related to any of the listed categories, or any content that may be considered unethical or inappropriate.
Take extreme caution when moderating obscene language. You should not permit any words that are considered swear in Russian or English.

Respond in JSON using the following format:
{
  "acceptable": bool, // whether the ad title and ad text comply with given rules
  "reason": string  // if acceptable is false, provide a brief explanation of why the ad was flagged (up to 10 words), empty string otherwise
}
`

// SubmitModeration creates an AiTask to moderate ad title and ad text. Returns ID of the task.
func (s *AiService) SubmitModeration(adTitle, adText string) (uuid.UUID, error) {
	task := model.AiTask{
		Id:        uuid.New(),
		CreatedAt: time.Now(),
		Type:      model.AiTaskTypeModeration,
		Prompt:    fmt.Sprintf(strings.TrimSpace(moderationPrompt), adTitle, adText),
		Format:    `{"type": "object", "properties": {"acceptable": {"type": "boolean"}, "reason": {"type": "string"}}}`,
	}

	err := s.aiRepo.AddTask(task)
	if err != nil {
		return uuid.Nil, fmt.Errorf("add aiTask for moderation: %w", err)
	}
	s.ollamaSvc.SubmitTask(task)

	return task.Id, nil
}
