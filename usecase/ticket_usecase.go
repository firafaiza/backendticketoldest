package usecase

import (
	"errors"
	"time"

	"ticket.narindo.com/model"
	"ticket.narindo.com/repository"
	"ticket.narindo.com/utils"
)

type TicketUseCase interface {
	CreateTicket(ticket *model.Ticket) (model.Ticket, error)
	ListAllTicket() ([]model.Ticket, error)
	ListByUser(userRole *model.UserRole) ([]model.Ticket, error)
	ListByDepartment(departmentId int) ([]model.Ticket, error)
	ListTicketSortedBy(userRole *model.UserRole, orderBy, field string) ([]model.Ticket, error)
	GetById(ticketId string) (model.Ticket, error)
	GetSummaryTicket(statusId int, userRole *model.UserRole) ([]model.Ticket, error)                                     //CHECK
	GetSummaryTicketByDate(statusId int, userRole *model.UserRole, startDate, endDate time.Time) ([]model.Ticket, error) // BELUM BISA
	UpdatePIC(ticketId string, picId int) error
	UpdateStatus(ticketId string, statusId int) error
}

type ticketUseCase struct {
	repo     repository.TicketRepository
	repoDept repository.PicDepartmentRepository
	repoPic  repository.PicRepository
}

func (t *ticketUseCase) CreateTicket(ticket *model.Ticket) (model.Ticket, error) {
	var toCheck = make(map[string]string)
	toCheck["TicketSubject"] = ticket.TicketSubject
	toCheck["TicketMessage"] = ticket.TicketMessage
	toCheck["CsId"] = ticket.CsId

	resultCheck := utils.CheckEmptyStringRequest(toCheck)
	if resultCheck != "" {
		return model.Ticket{}, errors.New("create failed, missing required field")
	}

	if ticket.CsId == "" || ticket.DepartmentId == 0 || ticket.TicketSubject == "" || ticket.TicketMessage == "" {
		return model.Ticket{}, errors.New("create failed, missing required field")
	}

	// if userRole.RoleId == 2 {
	// 	pic, _ := t.repoPic.FindBy(map[string]interface{}{"id": ticket.PicId})
	// 	if ticket.DepartmentId != pic.DepartmentId {
	// 		return model.Ticket{}, errors.New("create failed, must be in the same department id")
	// 	}
	// }

	department, err := t.repoDept.FindAllBy(map[string]interface{}{"id": ticket.DepartmentId})
	if err != nil {
		return model.Ticket{}, errors.New("problem exists")
	}
	if len(department) == 0 {
		return model.Ticket{}, errors.New("create failed, departement not found")
	}

	csId, err := t.repo.FindAllBy(map[string]interface{}{"cs_id": ticket.CsId})
	if err != nil {
		return model.Ticket{}, errors.New("problem exists")
	}
	if len(csId) == 0 {
		return model.Ticket{}, errors.New("create failed, cs id not found")
	}

	priority, err := t.repo.FindAllBy(map[string]interface{}{"priority_id": ticket.PriorityId})
	if err != nil {
		return model.Ticket{}, errors.New("problem exists")
	}
	if len(priority) == 0 {
		return model.Ticket{}, errors.New("create failed, priority not found")
	}

	var newTicketId string
	// diulang sampai ticketid baru gaada di db, kalo ada terus ngeloop sampai gaada yang kembar
	for {
		newTicketId = utils.GenerateId(department[0].DepartmentName)
		ticketIdExist, _ := t.repo.FindAllBy(map[string]interface{}{"ticket_id": newTicketId})
		if len(ticketIdExist) != 0 {
			newTicketId = utils.GenerateId(department[0].DepartmentName)
			continue
		}
		break
	}

	ticket.TicketId = newTicketId
	ticket.StatusId = 1
	ticket.IsActive = true

	err = t.repo.Create(ticket)

	return *ticket, err
}

func (t *ticketUseCase) ListAllTicket() ([]model.Ticket, error) {
	return t.repo.FindAll()
}

func (t *ticketUseCase) ListByUser(userRole *model.UserRole) ([]model.Ticket, error) {
	if userRole.RoleId == 2 {
		ticketList, err := t.repo.FindAllBy(map[string]interface{}{"cs_id": userRole.UserId})
		if err != nil {
			return []model.Ticket{}, err
		}
		return ticketList, nil
	}

	pic, _ := t.repoPic.FindAllBy(map[string]interface{}{"user_id": userRole.UserId})
	if len(pic) == 0 {
		return []model.Ticket{}, errors.New("can not find pic")
	}

	ticketList, err := t.repo.FindAllBy(map[string]interface{}{"pic_id": pic[0].Id})
	if err != nil {
		return []model.Ticket{}, err
	}
	return ticketList, nil
}

func (t *ticketUseCase) ListByDepartment(departmentId int) ([]model.Ticket, error) {
	ticketList, err := t.repo.FindAllBy(map[string]interface{}{"department_id": departmentId, "status_id": 1})
	if err != nil {
		return []model.Ticket{}, err
	}
	return ticketList, nil
}

func (t *ticketUseCase) ListTicketSortedBy(userRole *model.UserRole, orderBy, field string) ([]model.Ticket, error) {
	var tickets []model.Ticket
	if field == "created_at" || field == "updated_at" || field == "count_down" {
		if userRole.RoleId == 1 {
			return t.repo.SortBy(map[string]interface{}{}, map[string]interface{}{field: orderBy})
		} else if userRole.RoleId == 2 {
			return t.repo.SortBy(map[string]interface{}{"cs_id": userRole.UserId}, map[string]interface{}{field: orderBy})
		} //filterBy itu userid atau csid jika bukan superadmin, orderby itu asc atau desc, kalau field itu kolom mana yang mau diorder
		userPic, _ := t.repoPic.FindBy(map[string]interface{}{"user_id": userRole.UserId})
		return t.repo.SortBy(map[string]interface{}{"pic_id": userPic.DepartmentId}, map[string]interface{}{field: orderBy})
	}
	return tickets, nil
}

func (t *ticketUseCase) GetById(ticketId string) (model.Ticket, error) {
	return t.repo.FindBy(map[string]interface{}{"ticket_id": ticketId})
}

func (t *ticketUseCase) GetSummaryTicket(statusId int, userRole *model.UserRole) ([]model.Ticket, error) {
	if userRole.RoleId == 1 {
		return t.repo.FindAllBy(map[string]interface{}{"status_id": statusId})
	} else if userRole.RoleId == 2 {
		return t.repo.FindAllBy(map[string]interface{}{"status_id": statusId, "cs_id": userRole.UserId})
	}
	userPic, _ := t.repoPic.FindBy(map[string]interface{}{"user_id": userRole.UserId})
	return t.repo.FindAllBy(map[string]interface{}{"status_id": statusId, "pic_id": userPic.DepartmentId})
}

// belum di tes
func (t *ticketUseCase) GetSummaryTicketByDate(statusId int, userRole *model.UserRole, startDate, endDate time.Time) ([]model.Ticket, error) {
	if userRole.RoleId == 1 {
		whereClause := t.repo.WhereClauseByDate(nil, &statusId, nil, startDate, endDate)
		return t.repo.FindWithQuery(whereClause)
	} else if userRole.RoleId == 2 {
		whereClause := t.repo.WhereClauseByDate(&userRole.UserId, &statusId, nil, startDate, endDate)
		return t.repo.FindWithQuery(whereClause)
	}
	userPic, _ := t.repoPic.FindBy(map[string]interface{}{"user_id": userRole.UserId})
	whereClause := t.repo.WhereClauseByDate(&userRole.UserId, &statusId, &userPic.DepartmentId, startDate, endDate)
	return t.repo.FindWithQuery(whereClause)
}

func (t *ticketUseCase) UpdatePIC(ticketId string, picId int) error {
	return t.repo.UpdateBy(map[string]interface{}{"ticket_id": ticketId}, map[string]interface{}{"pic_id": picId})
}

func (t *ticketUseCase) UpdateStatus(ticketId string, statusId int) error {
	if statusId == 3 {
		ticket, _ := t.repo.FindBy(map[string]interface{}{"ticket_id": ticketId})
		var solveDuration = ticket.UpdatedAt.Day() - ticket.CreatedAt.Day()
		return t.repo.UpdateBy(map[string]interface{}{"ticket_id": ticketId}, map[string]interface{}{"status_id": statusId, "solving_duration": solveDuration})
	}
	return t.repo.UpdateBy(map[string]interface{}{"ticket_id": ticketId}, map[string]interface{}{"status_id": statusId})
}

func NewTicketUseCase(repo repository.TicketRepository, repoDept repository.PicDepartmentRepository, repoPic repository.PicRepository) TicketUseCase {
	return &ticketUseCase{
		repo:     repo,
		repoDept: repoDept,
		repoPic:  repoPic,
	}
}
