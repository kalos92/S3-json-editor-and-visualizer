<script setup>
import { ref } from 'vue'
import { NButton, NSwitch, NP, NInput, NSpace, NGrid, NGridItem, NCard, NDivider } from 'naive-ui'
import Modal from './Modal.vue';

const openErrorModal = ref(false)
const openSuccessModal = ref(false)
const bucket = ref('')
const S3Key = ref('')
const jsonBody = ref('')
const region = ref('')
const errorString = ref('')
const formatted_download = ref(false)
const athena = ref(false)
const successString = ref("File Uploaded!")
const Success = "Success"
const Error = "Reporting an error"
const previewFormat = ref(true)
const formatJson = ref(false)

function prettyPrintJson() {

  if (formatJson.value) {
    if (athena.value) {
      const values = jsonBody.value.split("\n")
      jsonBody.value = ""
      for (let i = 0; i < values.length; i++) {
        if (values[i] === "") {
          continue
        }
        console.log(values[i])
        const obj = JSON.parse(values[i]);

        jsonBody.value += JSON.stringify(obj, undefined, 2)
        jsonBody.value += "\n"
      }
      return
    }
    const obj = JSON.parse(jsonBody.value);
    jsonBody.value = ""
    jsonBody.value = JSON.stringify(obj, undefined, 2)
  } else {
    if (athena.value) {
      const replaced = jsonBody.value.replace("}\n{", "}${")
      const values = replaced.split("$")
      jsonBody.value = ""
      for (let i = 0; i < values.length; i++) {
        if (values[i] === "") {
          continue
        }
        console.log(values[i])
        const obj = JSON.parse(values[i]);
        jsonBody.value += JSON.stringify(obj)
        jsonBody.value += "\n"
      }
      return
    }
    const obj = JSON.parse(jsonBody.value);
    jsonBody.value = ""
    jsonBody.value = JSON.stringify(obj)
  }

}

function validate() {
  if (bucket.value !== '' && S3Key.value !== '' && S3Key.value.endsWith(".json")) {

    console.log(window.location.origin + "/api/request-json-to-s3?bucket=" + bucket.value + "&s3key=" + S3Key.value + "&region=" + region.value + "&athena=" + athena.value)

    fetch(window.location.origin + "/api/request-json-to-s3?bucket=" + bucket.value + "&s3key=" + S3Key.value + "&region=" + region.value + "&athena=" + athena.value)
      .then((res) => res.json())
      .then((data) => {

        if (data.error) {
          openErrorModal.value = true
          errorString.value = data.error_string
        }
        else {
          jsonBody.value = ""
          if (formatted_download.value) {
            jsonBody.value = JSON.stringify(data.response, undefined, 2);
          } else {
            for (let i = 0; i < data.response.length; i++) {
              jsonBody.value += JSON.stringify(data.response[i]);
              jsonBody.value += "\n"
            }
          }
          previewFormat.value = false

          let elmnt = document.getElementById('top');
          elmnt.scrollIntoView(false);
        }
      });


  } else {
    console.log("Invalid input")
  }
}

function save() {
  console.log(jsonBody.value)

  var myBody = {
    jsonBody: JSON.stringify(jsonBody.value),
    region: region.value,
    bucket: bucket.value,
    S3Key: S3Key.value,
    athena: athena.value,
    formatted: formatted_download.value
  }
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(myBody, undefined, 2)
  };

  console.log(requestOptions)
  fetch(window.location.origin + "/api/overwrite-json-to-s3", requestOptions)
    .then(response => response.json())
    .then((data) => {

      if (data.error) {
        openErrorModal.value = true
        errorString.value = data.error_string
      }
      else {
        openSuccessModal.value = true
        successString.value = data.response[0].message
      }
    }
    );
}

function closeErrorModal() {
  openErrorModal.value = !openErrorModal
}


function closeSuccessModal() {
  openSuccessModal.value = !openSuccessModal
}


</script>

<template>

  <n-grid cols="1 l:2 " item-responsive responsive="screen">


    <n-grid-item span="1 400:1 600:1 800:1">

      <n-card :bordered="true" class="fullscreen">
        <n-p v-if="bucket !== ''">Your bucket is : {{ bucket }}</n-p>
        <n-p v-else>Insert a bucket name</n-p>
        <n-input v-model:value="bucket" placeholder="bucket name" />

        <n-p v-if="S3Key !== ''">Requested Key is: {{ S3Key }}</n-p>
        <n-p v-else>Insert a S3 Key</n-p>
        <n-input v-model:value="S3Key" placeholder="S3 Key" />
        <n-p v-if="region !== ''">Requested region is: {{ region }}</n-p>
        <n-p v-else>Insert a valid AWS region code (eu-west-1, ap-southeast-1, ...)</n-p>
        <n-input v-model:value="region" placeholder="eu-west-1" />


        <n-card :bordered="false">
          <n-grid cols="2 l:2" item-responsive responsive="screen" x-gap="50">
            <n-grid-item span="1 400:1 600:1 800:1">
              <n-space justify="space-between" size="large">
                <n-p> Athena compliant Json</n-p>
                <n-switch v-model:value="athena" :disabled="formatted_download" />
              </n-space>
            </n-grid-item>

          </n-grid>

          <n-divider />
          <n-space justify="center">
            <n-button @click="validate">Load File</n-button>
          </n-space>
        </n-card>
      </n-card>

    </n-grid-item>

    <n-grid-item span="1 400:1 600:1 800:1">
      <n-card :bordered="true" class="fullscreen">
        <n-grid cols="1 l:1" item-responsive responsive="screen" x-gap="50">
          <n-grid-item span="1 400:1 600:1 800:1">
            <n-p style="white-space: pre-line;">Json Preview</n-p>

            <n-input v-model:value="jsonBody" type="textarea" class="area"
              placeholder="Here will be loaded the selected Json"></n-input>
          </n-grid-item>

          <n-grid-item span="1 400:1 600:1 800:1">

            <n-divider />
            <n-grid cols="2 l:2" item-responsive responsive="screen" x-gap="50">
              <n-grid-item span="1 400:1 600:1 800:1">
                <n-space justify="space-between" size="large">
                  <n-p> Work with a formatted Json</n-p>
                  <n-switch v-model:value="formatJson" :disabled="previewFormat" @update:value="prettyPrintJson" />
                </n-space>
              </n-grid-item>
              <n-grid-item span="1 400:1 600:1 800:1">
                <n-space justify="space-between" size="large">
                  <n-p> Upload formatted</n-p>
                  <n-switch v-model:value="formatted_download" :disabled="athena" />
                </n-space>
              </n-grid-item>
            </n-grid>

            <n-divider />



            <n-space justify="space-between">
              <n-button @click="save">Overwrite on S3</n-button>

            </n-space>
            <Modal :error="successString" :open="openSuccessModal" @close-modal="closeSuccessModal" :title="Success">
            </Modal>
            <Modal :error="errorString" :open="openErrorModal" @close-modal="closeErrorModal" :title="Error"></Modal>
          </n-grid-item>
        </n-grid>
      </n-card>
    </n-grid-item>
  </n-grid>
  <div id="top"></div>

</template>


<style scoped>
.area {
  width: 100%;
  min-height: 300px;
  max-height: 70vh;
}

.fullscreen {
  height: 90vh;
}
</style>