/*
This is where the implementation of automatically
generated resource configuration route goes
*/
package resourceConfig

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Globys031/PostgreScrutiniser/backend/cmd"
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
	Logger        *utils.Logger
	Configuration *cmd.Configuration
	// config        GetResourceConfigByIdParamsConfig
}

func (impl *ResourceConfigImpl) GetResourceConfigs(c *gin.Context) {
	// Reuse the same variable that contains resource setting details
	if impl.Configuration == nil {
		impl.Configuration = cmd.InitChecks(impl.Logger)
	}

	data := cmd.RunChecks(impl.Configuration, impl.Logger)
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusAccepted, "application/json", jsonData)
}
func (impl *ResourceConfigImpl) GetResourceConfigById(c *gin.Context, config GetResourceConfigByIdParamsConfig) {
	// Reuse the same variable that contains resource setting details
	if impl.Configuration == nil {
		impl.Configuration = cmd.InitChecks(impl.Logger)
	}

	configName := c.Param("config")
	var configData *cmd.ResourceSetting
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
