//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package repository

import (
	"time"

	"github.com/cloustone/pandas/apimachinery/models"
	"github.com/cloustone/pandas/pkg/cache"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	modelDB *gorm.DB
	cache   cache.Cache
}

// New return repository instance that manage device models and etc in dmms
func New(repositoryPath string, cache cache.Cache) *Repository {
	modelDB, err := gorm.Open(repositoryPath, "pandas-dmms.db")
	if err != nil {
		logrus.Fatal(err)
	}
	modelDB.AutoMigrate(&models.DeviceModel{})
	return &Repository{
		modelDB: modelDB,
		cache:   cache,
	}
}

// AddProject add a newly created project into repository
func (r *Repository) AddProject(principal *models.Principal, project *models.Project) (*models.Project, error) {
	db := r.modelDB.New()
	defer db.Close()
	project.CreatedAt = time.Now()
	project.LastUpdatedAt = time.Now()
	db.Save(project)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return project, nil
}

// GetProjects return user's project
func (r *Repository) GetProjects(principal *models.Principal, query *models.Query) ([]*models.Project, error) {
	db := r.modelDB.New()
	defer db.Close()
	projects := []*models.Project{}

	db.Where("userId = ?", principal.ID).Find(projects)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return projects, nil
}

// GetProject return specified project
func (r *Repository) GetProject(principal *models.Principal, projectID string) (*models.Project, error) {
	db := r.modelDB.New()
	defer db.Close()
	project := &models.Project{}

	db.Where("userId = ? AND projectId = ? ", principal.ID, projectID).Find(project)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return project, nil
}

// DeleteProject remove a project from repository
func (r *Repository) DeleteProject(principal *models.Principal, projectID string) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Delete(&models.Project{ProjectID: projectID})
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// UpdateProject update an already existed project
func (r *Repository) UpdateProject(principal *models.Principal, project *models.Project) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Where("userId = ? AND projectId = ?", principal.ID, project.ID).Find(nil)
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	project.LastUpdatedAt = time.Now()
	db.Save(project)
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// Devices
func (r *Repository) AddDevice(principal *models.Principal, device *models.Device) (*models.Device, error) {
	return nil, nil
}

func (r *Repository) DeleteDevice(principal *models.Principal, deviceID string, query *models.Query) error {
	return nil
}

func (r *Repository) GetDevices(principal *models.Principal, query *models.Query) ([]*models.Device, error) {
	return nil, nil
}

// View

// AddView add a newly created view into repository
func (r *Repository) AddView(principal *models.Principal, view *models.View) (*models.View, error) {
	db := r.modelDB.New()
	defer db.Close()
	view.CreatedAt = time.Now()
	view.LastUpdatedAt = time.Now()
	db.Save(view)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return view, nil
}

// GetViews return user's all view
func (r *Repository) GetViews(principal *models.Principal, query *models.Query) ([]*models.View, error) {
	db := r.modelDB.New()
	defer db.Close()
	views := []*models.View{}

	db.Where("userId = ?", principal.ID).Find(views)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return views, nil
}

// GetView return specified view
func (r *Repository) GetView(principal *models.Principal, viewID string) (*models.View, error) {
	db := r.modelDB.New()
	defer db.Close()
	view := &models.View{}

	db.Where("userId = ? AND viewId = ?", principal.ID, viewID).Find(view)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return view, nil
}

// DeleteView delete specified view
func (r *Repository) DeleteView(principal *models.Principal, viewID string) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Delete(&models.View{ViewID: viewID})
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// UpdateView update an already existed view
func (r *Repository) UpdateView(principal *models.Principal, view *models.View) (*models.View, error) {
	db := r.modelDB.New()
	defer db.Close()
	db.Where("userId = ? AND viewId = ?", principal.ID, view.ID).Find(nil)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	view.LastUpdatedAt = time.Now()
	db.Save(view)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return view, nil
}

// Workshop

// AddWorkshop add a newly created workshop into repository
func (r *Repository) AddWorkshop(principal *models.Principal, w *models.Workshop) (*models.Workshop, error) {
	db := r.modelDB.New()
	defer db.Close()
	w.CreatedAt = time.Now()
	w.LastUpdatedAt = time.Now()
	db.Save(w)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return w, nil
}

// GetViews return user's all workshops
func (r *Repository) GetWorkshops(principal *models.Principal, query *models.Query) ([]*models.Workshop, error) {
	db := r.modelDB.New()
	defer db.Close()
	ws := []*models.Workshop{}

	db.Where("userId = ?", principal.ID).Find(ws)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return ws, nil
}

// GetView return specified workshop
func (r *Repository) GetWorkshop(principal *models.Principal, wid string) (*models.Workshop, error) {
	db := r.modelDB.New()
	defer db.Close()
	w := &models.Workshop{}

	db.Where("userId = ? AND Id = ?", principal.ID, wid).Find(w)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return w, nil
}

// DeleteView delete specified workshop
func (r *Repository) DeleteWorkshop(principal *models.Principal, wid string) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Delete(&models.Workshop{WorkshopID: wid})
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// UpdateView update an already existed workshop
func (r *Repository) UpdateWorkshop(principal *models.Principal, w *models.Workshop) (*models.Workshop, error) {
	db := r.modelDB.New()
	defer db.Close()
	db.Where("userId = ? AND Id = ?", principal.ID, w.ID).Find(nil)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	w.LastUpdatedAt = time.Now()
	db.Save(w)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return w, nil
}

// Variable

// AddVariable add a newly created  variable into repository
func (r *Repository) AddVariable(principal *models.Principal, w *models.Variable) (*models.Variable, error) {
	db := r.modelDB.New()
	defer db.Close()
	w.CreatedAt = time.Now()
	w.LastUpdatedAt = time.Now()
	db.Save(w)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return w, nil
}

// GetViews return user's all variable
func (r *Repository) GetVariables(principal *models.Principal, query *models.Query) ([]*models.Variable, error) {
	db := r.modelDB.New()
	defer db.Close()
	variables := []*models.Variable{}

	db.Where("userId = ?", principal.ID).Find(variables)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return variables, nil
}

// GetView return specified variables
func (r *Repository) GetVariable(principal *models.Principal, wid string) (*models.Variable, error) {
	db := r.modelDB.New()
	defer db.Close()
	variable := &models.Variable{}

	db.Where("userId = ? AND Id = ?", principal.ID, wid).Find(variable)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return variable, nil
}

// DeleteView delete specified variable
func (r *Repository) DeleteVariable(principal *models.Principal, wid string) error {
	db := r.modelDB.New()
	defer db.Close()
	db.Delete(&models.Variable{ID: wid})
	if errs := db.GetErrors(); len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// UpdateVariable update an already existed variable
func (r *Repository) UpdateVariable(principal *models.Principal, variable *models.Variable) (*models.Variable, error) {
	db := r.modelDB.New()
	defer db.Close()
	db.Where("userId = ? AND Id = ?", principal.ID, variable.ID).Find(nil)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	variable.LastUpdatedAt = time.Now()
	db.Save(variable)
	if errs := db.GetErrors(); len(errs) > 0 {
		return nil, errs[0]
	}
	return variable, nil
}
