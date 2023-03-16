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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusAccepted, "application/json", jsonData)
}
func (impl *ResourceConfigImpl) GetResourceConfigById(c *gin.Context, config GetResourceConfigByIdParamsConfig) {
	// Reuse the same reference that contains resource setting details
	if impl.Configuration == nil {
		impl.Configuration = InitChecks(impl.DbHandler, impl.AppUser, impl.PostgresUser, impl.Logger)
	}

	configName := c.Param("config")
	var configData *ResourceSetting
	var err error
	switch configName {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No resource configuration with name: %s", configName)})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get suggestion. See /var/log/postgrescrutiniser/error.log for more details"})
		return
	}
	c.JSON(http.StatusAccepted, configData)
}

// Accepts a request body containing an array of suggestions
func (impl *ResourceConfigImpl) PatchResourceConfigs(c *gin.Context) {
	// Bind post body and validate
	suggestions := PatchResourceConfigsJSONBody{}
	if err := c.BindJSON(&suggestions); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "incorrect payload format"})
		return
	}
	if len(suggestions) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "empty payload array"})
		return
	}

	// Reuse the same reference that contains resource setting details
	if impl.Configuration == nil {
		impl.Configuration = InitChecks(impl.DbHandler, impl.AppUser, impl.PostgresUser, impl.Logger)
	}

	err := impl.Configuration.ApplySuggestions(&suggestions, impl.Logger)
	if err != nil {
		errorMessage := fmt.Sprintf("%s. See /var/log/postgrescrutiniser/error.log for more details", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}
}
