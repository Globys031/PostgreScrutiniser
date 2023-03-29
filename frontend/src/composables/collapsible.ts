// Function necessary for smooth collapse transition
export const toggleCollapseContent = (content: HTMLElement | null) => {
  if (content) {
    const scrollHeight = content.scrollHeight;
    const maxHeight = content.style.maxHeight;

    maxHeight === "" || maxHeight === "0px"
      ? (content.style.maxHeight = scrollHeight + "px")
      : (content.style.maxHeight = "0");
  } else {
    console.error("content ref is undefined");
  }
};
