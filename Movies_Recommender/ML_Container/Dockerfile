FROM jupyter/datascience-notebook


# Setting the Work Directory
WORKDIR /app

# Copy the "model.ipynb" into the Container
COPY model.ipynb ./

# Exposing Port 8888 
EXPOSE 8888

# Starts Jupyter Notebook with certain arguments.
CMD ["jupyter", "notebook", "--ip='*'", "--port=8888", "--no-browser", "--allow-root"]

