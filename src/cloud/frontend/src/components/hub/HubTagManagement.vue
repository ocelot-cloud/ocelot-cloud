<template>
  <div class="container mt-5 col-lg-6 col-md-8 col-sm-10">
    <div class="entity-management-container p-4 shadow-sm bg-dark rounded">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h4>Tag Management</h4>
        <button id="button-back-to-app" class="btn btn-secondary" @click="goToHubPage('/')">Back to App Management</button>
      </div>

      <div class="app-info mb-3">
        <p class="mb-1">App: <strong>{{ app }}</strong></p>
      </div>

      <div class="file-upload-area mb-4">
        <input type="file" ref="fileInput" @change="handleFileUpload" class="d-none" />
        <div
            id="drag-and-drop-area"
            class="drop-zone d-flex justify-content-center align-items-center bg-secondary bg-opacity-25 rounded"
            @dragover.prevent
            @drop.prevent="handleDrop"
        >
          <p class="m-0">Drag and drop your file here</p>
        </div>
      </div>
      <div v-if="submitted" class="invalid-feedback d-block mb-4">
        {{ errorMessageText }}
      </div>

      <div class="tag-list-section mb-4">
        <h4>Tag List</h4>
        <p v-if="!tagList || tagList.length === 0">(No tags created yet)</p>

        <div class="d-flex justify-content-center">
          <ul id="tag-list" class="list-group w-100">
            <li
                v-for="tag in tagList"
                :key="tag"
                class="list-group-item d-flex justify-content-between align-items-center bg-secondary bg-opacity-25 text-white"
                :class="{ active: selectedTag === tag }"
                @click="selectTag(tag)"
                style="cursor: pointer;"
            >
              <span>{{ tag }}</span>
              <i v-if="selectedTag === tag" class="bi bi-check-circle-fill text-success"></i>
            </li>
          </ul>
        </div>
      </div>

      <div v-if="tagList && selectedTag" class="app-operations d-flex justify-content-end">
        <button id="button-download-tag" @click="downloadTag" class="btn btn-primary me-2">Download</button>
        <button id="button-delete-tag" @click="showDeleteConfirmation = true" class="btn btn-danger">Delete</button>
      </div>

      <HubDeletionConfirmationDialog
          v-model:visible="showDeleteConfirmation"
          :on-confirm="deleteTag"
          messageSuffix="this tag?"
      />
    </div>
  </div>
</template>


<script lang="ts">
import {defineComponent, onMounted, ref} from 'vue';
import { useRoute } from 'vue-router';
import {
  alertError,
  defaultAllowedSymbols,
  doRequest, generateInvalidInputMessage, getDefaultValidationRegex,
  goToHubPage,
  defaultMaxLength,
  defaultMinLength, tagAllowedSymbols
} from "@/components/hub/shared";
import HubDeletionConfirmationDialog from "@/components/hub/HubDeletionConfirmationDialog.vue";

export default defineComponent({
  name: "HubTagManagement",
  components: {HubDeletionConfirmationDialog},
  methods: {goToHubPage},

  setup() {
    const tagList = ref<string[]>([]);
    const route = useRoute();
    const app = route.query.app
    const user = route.query.user
    const selectedTag = ref("");
    const showDeleteConfirmation = ref(false);
    const submitted = ref(false);
    const errorMessageText = ref(generateInvalidInputMessage("tag", tagAllowedSymbols, defaultMinLength, defaultMaxLength))

    const handleFileUpload = (event: Event) => {
      const files = (event.target as HTMLInputElement).files;
      if (files && files.length > 0) {
        uploadFile(files[0]);
      }
    };

    const uploadFile = (file: File) => {
      const suffix = '.tar.gz';
      if (!file.name.endsWith(suffix)) {
        alert(`The file must have a ${suffix} suffix.`);
        return;
      }

      const tag = file.name.slice(0, -suffix.length);

      let regex = new RegExp(`^${tagAllowedSymbols}{${defaultMinLength},${defaultMaxLength}}$`)
      if (!regex.test(tag)) {
        submitted.value = true
        return;
      }

      const reader = new FileReader();
      reader.onload = async (event) => {
        const content = btoa(
            String.fromCharCode(...new Uint8Array(event.target?.result as ArrayBuffer))
        );
        const tagUpload = {app, tag, content};

        const response = await doRequest("/tags/upload", tagUpload)
        if (response) {
          submitted.value = false
          getTags()
        }

      };

      reader.onerror = () => {
        console.error('Error reading file');
      };
      reader.readAsArrayBuffer(file);
    };

    const getTags = async () => {
      const response = await doRequest("/tags/get-tags", { user, app })
      if (response != null) {
        tagList.value = response.data as string[];
        if (tagList.value != null) {
          tagList.value.sort()
        }
      }
    };

    const deleteTag = async () => {
      await doRequest("/tags/delete", {app, tag: selectedTag.value})
      await getTags()
      showDeleteConfirmation.value = false
    }

    const downloadTag = async () => {
      try {
        const response = await doRequest("/tags/", { user, app, tag: selectedTag.value })
        if (response != null) {
          const blob = new Blob([response.data], { type: 'application/gzip' });
          const downloadUrl = window.URL.createObjectURL(blob);
          const link = document.createElement('a');
          link.href = downloadUrl;
          link.setAttribute('download', `${selectedTag.value}.tar.gz`);
          document.body.appendChild(link);
          link.click();
          link.remove();
          console.log("File download started successfully");
        }
      } catch (error) {
        alertError(error)
        console.error('Error during file download:', error);
      }
    };

    const selectTag = (tag: string) => {
      if (selectedTag.value == tag) {
        selectedTag.value = ""
      } else {
        selectedTag.value = tag;
      }
    }

    const handleDrop = (event: DragEvent) => {
      const files = event.dataTransfer?.files;
      if (files && files.length > 0) {
        uploadFile(files[0]);
      }
    };

    onMounted(() => {
      getTags()
    });

    return {
      handleFileUpload,
      tagList,
      app,
      selectedTag,
      deleteTag,
      selectTag,
      downloadTag,
      showDeleteConfirmation,
      handleDrop,
      errorMessageText,
      submitted
    }
  },
});
</script>
