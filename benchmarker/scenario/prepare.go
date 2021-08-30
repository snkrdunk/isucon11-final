package scenario

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/isucon/isucon11-final/benchmarker/api"

	"github.com/isucon/isucandar/parallel"

	"github.com/isucon/isucandar/worker"

	"github.com/isucon/isucon11-final/benchmarker/generate"

	"github.com/isucon/isucon11-final/benchmarker/model"

	"github.com/isucon/isucandar/agent"
	"github.com/isucon/isucandar/failure"
	"github.com/isucon/isucon11-final/benchmarker/fails"

	"github.com/isucon/isucandar"
)

const (
	prepareTimeout           = 20
	AnnouncementCountPerPage = 20
)

func (s *Scenario) Prepare(ctx context.Context, step *isucandar.BenchmarkStep) error {
	ContestantLogger.Printf("===> PREPARE")

	a, err := agent.NewAgent(
		agent.WithNoCache(),
		agent.WithNoCookie(),
		agent.WithTimeout(20*time.Second),
		agent.WithBaseURL(s.BaseURL.String()),
	)
	if err != nil {
		return failure.NewError(fails.ErrCritical, err)
	}

	a.Name = "benchmarker-initializer"

	ContestantLogger.Printf("start Initialize")
	_, err = InitializeAction(ctx, a)
	if err != nil {
		ContestantLogger.Printf("initializeが失敗しました")
		return failure.NewError(fails.ErrCritical, err)
	}

	err = s.prepareNormal(ctx, step)
	if err != nil {
		return err
	}

	errors := step.Result().Errors
	hasErrors := func() bool {
		errors.Wait()
		return len(errors.All()) > 0
	}

	if hasErrors() {
		step.AddError(failure.NewError(fails.ErrCritical, fmt.Errorf("アプリケーション互換性チェックに失敗しました")))
		return nil
	}

	return nil
}

func (s *Scenario) prepareCheck(parent context.Context, step *isucandar.BenchmarkStep) error {
	initializeAgent, err := agent.NewAgent(
		agent.WithNoCache(),
		agent.WithNoCookie(),
		agent.WithTimeout(20*time.Second),
		agent.WithBaseURL(s.BaseURL.String()),
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	_, err = InitializeAction(ctx, initializeAgent)
	if err != nil {
		return err
	}

	//studentAgent, err := agent.NewAgent(agent.WithTimeout(prepareTimeout))
	//if err != nil {
	//	return err
	//}
	//student := s.prepareNewStudent()
	//student.Agent = studentAgent
	//
	//teacherAgent, err := agent.NewAgent(agent.WithTimeout(prepareTimeout))
	//if err != nil {
	//	return err
	//}
	//teacher := s.prepareNewTeacher()
	//teacher.Agent = teacherAgent

	return nil
}

func (s *Scenario) prepareNormal(ctx context.Context, step *isucandar.BenchmarkStep) error {
	const (
		prepareTeacherCount        = 2
		prepareStudentCount        = 2
		prepareCourseCount         = 20
		prepareCourseRegisterLimit = 20
		prepareClassCountPerCourse = 5
	)
	errors := step.Result().Errors
	hasErrors := func() bool {
		errors.Wait()
		return len(errors.All()) > 0
	}

	teachers := make([]*model.Teacher, 0, prepareTeacherCount)
	// TODO: ランダムなので同じ教師が入る可能性がある
	for i := 0; i < prepareTeacherCount; i++ {
		teachers = append(teachers, s.GetRandomTeacher())
	}

	students := make([]*model.Student, 0, prepareStudentCount)
	for i := 0; i < prepareStudentCount; i++ {
		userData, err := s.studentPool.newUserData()
		if err != nil {
			return err
		}
		students = append(students, model.NewStudent(userData, s.BaseURL, prepareCourseRegisterLimit))
	}

	courses := make([]*model.Course, 0, prepareCourseCount)
	mu := sync.Mutex{}
	// 教師のログインとコース登録をするワーカー
	w, err := worker.NewWorker(func(ctx context.Context, i int) {
		teacher := teachers[i%len(teachers)]
		_, err := LoginAction(ctx, teacher.Agent, teacher.UserAccount)
		if err != nil {
			AdminLogger.Printf("teacherのログインに失敗しました")
			step.AddError(failure.NewError(fails.ErrCritical, err))
			return
		}

		_, getMeRes, err := GetMeAction(ctx, teacher.Agent)
		if err != nil {
			AdminLogger.Printf("teacherのユーザ情報取得に失敗しました")
			step.AddError(err)
			return
		}
		if err := verifyMe(&getMeRes, teacher.UserAccount, true); err != nil {
			step.AddError(err)
			return
		}

		param := generate.CourseParam(teacher, generate.WithPeriod(i%6), generate.WithDayOfWeek((i/6)+1))
		_, res, err := AddCourseAction(ctx, teacher, param)
		if err != nil {
			step.AddError(err)
			return
		}
		course := model.NewCourse(param, res.ID, teacher)
		mu.Lock()
		courses = append(courses, course)
		mu.Unlock()
	}, worker.WithLoopCount(prepareCourseCount))

	if err != nil {
		step.AddError(err)
		return err
	}

	w.Process(ctx)
	w.Wait()

	if hasErrors() {
		return failure.NewError(fails.ErrCritical, fmt.Errorf("アプリケーション互換性チェックに失敗しました"))
	}

	// 生徒のログインとコース登録
	w, err = worker.NewWorker(func(ctx context.Context, i int) {
		student := students[i]
		_, err := LoginAction(ctx, student.Agent, student.UserAccount)
		if err != nil {
			return
		}
		if err != nil {
			AdminLogger.Printf("studentのログインに失敗しました")
			step.AddError(failure.NewError(fails.ErrCritical, err))
			return
		}

		_, getMeRes, err := GetMeAction(ctx, student.Agent)
		if err != nil {
			AdminLogger.Printf("studentのユーザ情報取得に失敗しました")
			step.AddError(err)
			return
		}
		if err := verifyMe(&getMeRes, student.UserAccount, false); err != nil {
			step.AddError(err)
			return
		}

		_, err = TakeCoursesAction(ctx, student.Agent, courses)
		if err != nil {
			step.AddError(err)
			return
		}
		for _, course := range courses {
			student.AddCourse(course)
			course.AddStudent(student)
		}
	}, worker.WithLoopCount(prepareStudentCount))
	if err != nil {
		step.AddError(err)
		return err
	}
	w.Process(ctx)
	w.Wait()
	if hasErrors() {
		return failure.NewError(fails.ErrCritical, fmt.Errorf("アプリケーション互換性チェックに失敗しました"))
	}

	// コースのステータスの変更
	w, err = worker.NewWorker(func(ctx context.Context, i int) {
		course := courses[i]
		teacher := course.Teacher()
		_, err = SetCourseStatusInProgressAction(ctx, teacher.Agent, course.ID)
		if err != nil {
			step.AddError(err)
			return
		}
	}, worker.WithLoopCount(prepareCourseCount))
	if err != nil {
		step.AddError(err)
		return err
	}
	w.Process(ctx)
	w.Wait()

	if hasErrors() {
		return failure.NewError(fails.ErrCritical, fmt.Errorf("アプリケーション互換性チェックに失敗しました"))
	}

	// クラス追加、課題提出、ダウンロード、採点（お知らせは見ない）
	// workerはコースごと
	w, err = worker.NewWorker(func(ctx context.Context, i int) {
		course := courses[i]
		teacher := course.Teacher()

		for classPart := 0; classPart < prepareClassCountPerCourse; classPart++ {
			classParam := generate.ClassParam(course, uint8(classPart+1))
			_, classRes, err := AddClassAction(ctx, teacher.Agent, course, classParam)
			if err != nil {
				step.AddError(err)
				return
			}
			class := model.NewClass(classRes.ClassID, classParam)
			course.AddClass(class)

			announcement := generate.Announcement(course, class)
			_, ancRes, err := SendAnnouncementAction(ctx, teacher.Agent, announcement)
			if err != nil {
				step.AddError(err)
				return
			}
			announcement.ID = ancRes.ID
			course.BroadCastAnnouncement(announcement)

			courseStudents := course.Students()

			// 課題提出,
			// 生徒ごとのループ
			p := parallel.NewParallel(ctx, prepareStudentCount)
			for _, student := range courseStudents {
				student := student
				p.Do(func(ctx context.Context) {
					submissionData, fileName := generate.SubmissionData(course, class, student.UserAccount)
					_, err := SubmitAssignmentAction(ctx, student.Agent, course.ID, class.ID, fileName, submissionData)
					if err != nil {
						step.AddError(err)
						fmt.Println(class.ID, student.Code)
						return
					}
					submissionSummary := model.NewSubmission(fileName, submissionData, true)
					class.AddSubmission(student.Code, submissionSummary)
				})
			}
			p.Wait()

			// 課題ダウンロード
			_, assignmentsData, err := DownloadSubmissionsAction(ctx, teacher.Agent, course.ID, class.ID)
			if err != nil {
				step.AddError(err)
				return
			}
			if err := verifyAssignments(assignmentsData, class); err != nil {
				step.AddError(err)
				return
			}

			// 採点
			scores := make([]StudentScore, 0, len(students))
			for _, student := range students {
				sub := class.GetSubmissionByStudentCode(student.Code)
				if sub == nil {
					step.AddError(failure.NewError(fails.ErrCritical, fmt.Errorf("cannot find submission")))
					return
				}
				score := rand.Intn(101)
				sub.SetScore(score)
				scores = append(scores, StudentScore{
					score: score,
					code:  student.Code,
				})
			}

			_, err = PostGradeAction(ctx, teacher.Agent, course.ID, class.ID, scores)
			if err != nil {
				step.AddError(err)
				return
			}
		}
	}, worker.WithLoopCount(prepareCourseCount))
	if err != nil {
		step.AddError(err)
		return err
	}
	w.Process(ctx)
	w.Wait()

	if hasErrors() {
		return failure.NewError(fails.ErrCritical, fmt.Errorf("アプリケーション互換性チェックに失敗しました"))
	}

	w, err = worker.NewWorker(func(ctx context.Context, i int) {
		student := students[i]
		expected := student.Announcements()

		// createdAtが新しい方が先頭に来るようにソート
		sort.Slice(expected, func(i, j int) bool {
			return expected[i].Announcement.CreatedAt > expected[j].Announcement.CreatedAt
		})
		_, err := prepareCheckAnnouncementsList(ctx, student.Agent, "", expected, len(expected))
		if err != nil {
			step.AddError(err)
			return
		}

	}, worker.WithLoopCount(prepareStudentCount))
	if err != nil {
		step.AddError(err)
		return err
	}
	w.Process(ctx)
	w.Wait()
	if hasErrors() {
		return failure.NewError(fails.ErrCritical, fmt.Errorf("アプリケーション互換性チェックに失敗しました"))
	}

	return nil
}

func prepareCheckAnnouncementsList(ctx context.Context, a *agent.Agent, path string, expected []*model.AnnouncementStatus, expectedUnreadCount int) (prev string, err error) {
	errHttp := failure.NewError(fails.ErrCritical, fmt.Errorf("/api/announcements へのリクエストが失敗しました"))

	hres, res, err := GetAnnouncementListAction(ctx, a, path)
	if err != nil {
		return "", errHttp
	}
	prev, next := parseLinkHeader(hres)

	// 次のページが存在しない
	if len(res.Announcements) < AnnouncementCountPerPage || len(expected) < AnnouncementCountPerPage || next == "" {
		err = prepareCheckAnnouncementContent(expected, res, expectedUnreadCount)
		if err != nil {
			return "", err
		}
		return prev, nil
	}

	err = prepareCheckAnnouncementContent(expected[:AnnouncementCountPerPage], res, expectedUnreadCount)
	if err != nil {
		return "", err
	}

	_prev, err := prepareCheckAnnouncementsList(ctx, a, next, expected[AnnouncementCountPerPage:], expectedUnreadCount)
	if err != nil {
		return "", err
	}

	hres, res, err = GetAnnouncementListAction(ctx, a, _prev)
	if err != nil {
		return "", errHttp
	}

	err = prepareCheckAnnouncementContent(expected[:AnnouncementCountPerPage], res, expectedUnreadCount)
	if err != nil {
		return "", err
	}

	return prev, nil
}

func prepareCheckAnnouncementContent(expected []*model.AnnouncementStatus, actual api.GetAnnouncementsResponse, expectedUnreadCount int) error {
	errNotSorted := failure.NewError(fails.ErrCritical, fmt.Errorf("/api/announcements の順序が不正です"))
	//errNotMatchOver := failure.NewError(fails.ErrCritical, fmt.Errorf("存在しないはずの Announcement が見つかりました"))
	errNotMatch := failure.NewError(fails.ErrCritical, fmt.Errorf("announcement が期待したものと一致しませんでした"))
	//errHttp := failure.NewError(fails.ErrCritical, fmt.Errorf("/api/announcements へのリクエストが失敗しました"))
	errNoCount := failure.NewError(fails.ErrCritical, fmt.Errorf("announcement の数が期待したものと一致しませんでした"))

	if len(expected) != len(actual.Announcements) {
		return errNoCount
	}

	if expected == nil && actual.Announcements == nil {
		return nil
	} else if (expected == nil && actual.Announcements != nil) || (expected != nil && actual.Announcements == nil) {
		return errNotMatch
	}

	lastCreatedAt := int64(math.MaxInt64)
	for _, announcement := range actual.Announcements {
		// 順序の検証
		if lastCreatedAt < announcement.CreatedAt {
			return errNotSorted
		}
		lastCreatedAt = announcement.CreatedAt
	}
	for i := 0; i < len(actual.Announcements); i++ {
		expect := expected[i].Announcement
		actual := actual.Announcements[i]
		if !AssertEqual("announcement unread", expected[i].Unread, actual.Unread) ||
			!AssertEqual("announcement ID", expect.ID, actual.ID) ||
			!AssertEqual("announcement Code", expect.CourseID, actual.CourseID) ||
			!AssertEqual("announcement Title", expect.Title, actual.Title) ||
			!AssertEqual("announcement CourseName", expect.CourseName, actual.CourseName) ||
			!AssertEqual("announcement CreatedAt", expect.CreatedAt, actual.CreatedAt) {
			AdminLogger.Printf("extra announcements ->name: %v, title:  %v", actual.CourseName, actual.Title)
			return errNotMatch
		}
	}

	return nil
}
