<template>
  <h4>Tag Management</h4>
  <button id="button-back-to-app" class="btn btn-secondary" @click="goToHubPage('/')">Back to App Management</button>

  <p>App is: {{ app }}</p>

  <div>
    <input type="file" ref="fileInput" @change="handleFileUpload" style="display: none;" />
    <div
        id="drag-and-drop-area"
        class="drop-zone"
        @dragover.prevent
        @drop.prevent="handleDrop"
    >
      Drag and drop your file here
    </div>
  </div>

  <h4>Tag List</h4>
  <p v-if="tagList == null || tagList.length == 0"> (no apps created yet) </p>

  <div class="d-flex justify-content-center">
    <ul id="tag-list" class="list-group">
      <li
          v-for="(tag, index) in tagList"
          :key="tag"
          class="list-group-item"
          :class="{ active: selectedTag === tag }"
          @click="selectTag(tag)"
          style="cursor: pointer;"
      >
        {{ index + 1 }}) {{ tag }}
      </li>
    </ul>
  </div>
  <br>
  <div v-if="tagList != null && selectedTag != ''">
    <h4>App Operations</h4>
    <!-- TODO There should be a confirmation dialog previously -->
    <button id="button-delete-tag" @click="deleteTag" class="btn btn-danger ms-2">Delete</button>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, ref} from 'vue';
import axios from "axios";
import { useRoute } from 'vue-router';
import {goToHubPage} from "@/components/hub/shared";

export default defineComponent({
  name: "HubTagManagement",
  methods: {goToHubPage},

  setup() {
    const tagList = ref<string[]>([]);
    const route = useRoute();
    const app = route.query.app
    const user = route.query.user
    const selectedTag = ref<string>("");

    const handleFileUpload = (event: Event) => {
      const files = (event.target as HTMLInputElement).files;
      if (files && files.length > 0) {
        uploadFile(files[0]);
      }
    };

    // TODO Is that necessary?
    const handleDrop = (event: DragEvent) => {
      const files = event.dataTransfer?.files;
      if (files && files.length > 0) {
        uploadFile(files[0]);
      }
    };

    // TODO There should be bright styling when hover the drag and drop area with a file
    const uploadFile = (file: File) => {
      const suffix = '.tar.gz';

      if (!file.name.endsWith(suffix)) {
        alert(`The file must have a ${suffix} suffix.`);
        return;
      }
      const tag = file.name.slice(0, -suffix.length);

      const reader = new FileReader();
      reader.onload = async (event) => {
        const content = btoa(
            String.fromCharCode(...new Uint8Array(event.target?.result as ArrayBuffer))
        );
        const tagUpload = {app, tag, content};

        try {
          const response = await axios.post('http://localhost:8082/tags', tagUpload);
          if (response.status === 200) {
            console.log('File uploaded successfully');
          }
          await getTags()
        } catch (error) {
          console.error('Error uploading file:', error);
        }
      };

      reader.onerror = () => {
        console.error('Error reading file');
      };

      reader.readAsArrayBuffer(file); // Trigger reading the file as an ArrayBuffer
    };

    const getTags = async () => {
      const url = 'http://localhost:8082';
      try {
        let userAndApp = { user: user, app: app } // TODO Can be shortened I guess
        const response = await axios.post(url + '/tags/get-tags', userAndApp);
        if (response.status === 200) {
          tagList.value = response.data as string[];
          console.log("received tags: ", tagList.value)
        }
      } catch (error) {
        console.log("todo")
      }
    };

    const deleteTag = () => {
      console.log("todo")
    }

    const selectTag = (tag: string) => {
      if (selectedTag.value == tag) {
        selectedTag.value = ""
      } else {
        selectedTag.value = tag;
      }
    }

    onMounted(() => {
      getTags()
    });

    return {
      handleFileUpload,
      handleDrop,
      tagList,
      app,
      selectedTag,
      deleteTag,
      selectTag,
    }
  },
});
</script>