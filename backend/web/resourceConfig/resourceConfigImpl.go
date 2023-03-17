/*
This is where the implementation of automatically
generated resource configuration route goes
*/
package resourceConfig

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Globys031/PostgreScrutiniser/backend/utils"
	"github.com/gin-gonic/gin"
	// "github.com/Globys031/plotzemis/go/auth"
	// "github.com/Globys031/plotzemis/go/db"
)

/*
After `oapi-codegen` was used to generate API endpoints,
`ServerInterface` can be used to add our actual implementation
*/
type ResourceConfigImpl struct {
	AppUser       *utils.User
	PostgresUser  *utils.User
	Logger        *utils.Logger
	Configuration *Configuration
	DbHandler     *sql.DB
}

func (impl *ResourceConfigImpl) GetResourceConfigs(c *gin.Context) {
	// Reuse the same variable that contains resource setting details
	if impl.Configuration == nil {
		impl.Configuration = InitChecks(impl.DbHandler, impl.AppUser, impl.PostgresUser, impl.Logger)
	}

	data := RunChecks(impl.Configuration, impl.Logger)
	jsonData, err := json.Marshal(data)
	if err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, &errorMsg)
		return
	}

	c.Data(http.StatusAccepted, "application/json", jsonData)
}

func (impl *ResourceConfigImpl) GetResourceConfigById(c *gin.Context, config GetResourceConfigByIdParamsConfig) {
	// Reuse the same reference that contains resource setting details
	if impl.Configuration == nil {
		impl.Configuration = InitChecks(impl.DbHandler, impl.AppUser, impl.PostgresUser, impl.Logger)
	}

	var configData *ResourceSetting
	var err error
	switch config {
	case "autovacuum_work_mem":
		configData, err = impl.Configuration.CheckAutovacuumWorkMem(impl.Logger)
	case "dynamic_shared_memory_type":
		configData, err = impl.Configuration.CheckDynamicSharedMemoryType(impl.Logger)
	case "hash_mem_multiplier":
		configData, err = impl.Configuration.CheckHashMemMultiplier(impl.Logger)
	case "huge_page_size":
		configData, err = impl.Configuration.CheckHugePageSize(impl.Logger)
	case "huge_pages":
		configData, err = impl.Configuration.CheckHugePages(impl.Logger)
	case "maintenance_work_mem":
		configData, err = impl.Configuration.CheckMaintenanceWorkMem(impl.Logger)
	case "max_prepared_transactions":
		configData, err = impl.Configuration.CheckMaxPreparedTransactions(impl.Logger)
	case "max_stack_depth":
		configData, err = impl.Configuration.CheckMaxStackDepth(impl.Logger)
	case "shared_buffers":
		configData, err = impl.Configuration.CheckSharedBuffers(impl.Logger)
	case "shared_memory_type":
		configData, err = impl.Configuration.CheckSharedMemoryType(impl.Logger)
	case "temp_buffers":
		configData, err = impl.Configuration.CheckTempBuffers(impl.Logger)
	case "work_mem":
		configData, err = impl.Configuration.CheckWorkMem(impl.Logger)
	case "logical_decoding_work_mem":
		configData, err = impl.Configuration.ChecklogicalDecodingWorkMem(impl.Logger)
	default:
		errorMsg := &ErrorMessage{
			ErrorMessage: fmt.Sprintf("No resource configuration with name: %s", config),
		}
		c.JSON(http.StatusBadRequest, errorMsg)
		return
	}

	if err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: "Could not get suggestion. See /var/log/postgrescrutiniser/error.log for more details",
		}
		c.JSON(http.StatusInternalServerError, errorMsg)
		return
	}
	c.JSON(http.StatusAccepted, configData)
}

// Accepts a request body containing an array of suggestions
func (impl *ResourceConfigImpl) PatchResourceConfigs(c *gin.Context) {
	// Bind post body and validate
	suggestions := PatchResourceConfigsJSONBody{}
	if err := c.BindJSON(&suggestions); err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: "incorrect payload format",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorMsg)
		return
	}
	if len(suggestions) == 0 {
		errorMsg := &ErrorMessage{
			ErrorMessage: "empty payload array",
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errorMsg)
		return
	}

	// Reuse the same reference that contains resource setting details
	if impl.Configuration == nil {
		impl.Configuration = InitChecks(impl.DbHandler, impl.AppUser, impl.PostgresUser, impl.Logger)
	}

	err := impl.Configuration.ApplySuggestions(&suggestions, impl.Logger)
	if err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: fmt.Sprintf("%s. See /var/log/postgrescrutiniser/error.log for more details", err.Error()),
		}
		c.JSON(http.StatusInternalServerError, errorMsg)
		return
	}
}

func (impl *ResourceConfigImpl) DeleteResourceConfigs(c *gin.Context) {
	if err := impl.Configuration.DiscardConfigs(impl.Logger); err != nil {
		errorMsg := &ErrorMessage{
			ErrorMessage: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, errorMsg)
	}
}
