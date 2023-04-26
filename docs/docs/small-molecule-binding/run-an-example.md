---
title: Run an example
sidebar_position: 4
---

import AsciinemaPlayer from '../../src/components/AsciinemaPlayer.js';

This tutorial will guide you through the steps required to run a small molecule binding tool directly from your computer.

By the end of this tutorial, you will have:

* Run a docking tool (Equibind) on our provided test data, which include a protein file and a small molecule file
* Visualized the results

Let’s get started!

---

**Time needed:**
- 5 minutes

**Requirements:**

- Install PLEX ([installation guide here](../getting-started/install-plex.md))
- No previous technical experience - we’ll walk through each step.

---

After [installing PLEX](../getting-started/install-plex.md), follow the steps below:

### 1. Submit a job

Copy and paste the following command to run the tool using our provided test data and press **Enter**:

```
./plex -tool equibind -input-dir ./testdata/binding/abl/
```

You should see the following:

<div style={{ marginBottom: '20px' }}>
    <AsciinemaPlayer 
        src="/terminal-recordings/small-molecule-binding-run-an-example.cast"
        rows={20}
        idleTimeLimit={3}
        preload={true}
        autoPlay={true}
        loop={false}
        speed={2.5}
    />
</div>

:::tip

You might get a pop-up asking *"Do you want the application “plex” to accept incoming network connections?”*. Click ***“Allow”***.

If you need to, you can turn off your firewall. To do this on your Mac, go to settings via `System Preferences > Security & Privacy`. Then go to the Firewall tab, and click the padlock icon at the bottom of the window to make changes. Click the “Turn off Firewall” button and try running the tool again.

:::

### 2. Get the results

Once the job is complete and the results have downloaded, you will see the file path where your results can be found. It will look something like this: 

![result.png](screenshot-resultsdownloaded.png)

To open the folder where your results are stored, type ```open ``` into your command line, followed by the file path you were given as an output e.g. ```open /Users/user-demo-account/plex/bd5c4751-0a7a-42bd-a92d-6a1a0758a6a3```.  Press **Enter**

This will show your results in Finder.

![InstallationTutorial_Screenshot_with_results_folder](InstallationTutorial_Screenshot_with_results_folder.png)

### 3. Visualize the job results

To visualize the results, we are going to use Molstar.

In your results Finder window, click the “combined_results” folder, then the “outputs” folder.

(In our example, the file path would be: ```/Users/user-demo-account/plex/8deeb1b6-d53f-44e4-8e0a-6d7be6f1c43d/combined_results/outputs)```

You should see:

![InstallationTutorial_Screenshot_with_outputs](InstallationTutorial_Screenshot_with_outputs.png)

To inspect the results interactively with a viewer, open [the Molstar visualizer in your web browser.](https://molstar.org/viewer/)

Drag and drop **both the files** into the central blank frame in Molstar to see the result as per the gif below:

![molstar_draganddrop](Gif_-_drag_and_drop_molstar.gif)

Here is a close up what the result looks like in Molstar:

![molstar](InstallationTutorial_Screenshot_of_Molstar.png)

You can see how tightly the small molecule is predicted to bind to the protein. 

For more on how to use the Molstar viewer, check out [the Molstar documentation](https://molstar.org/viewer-docs/).

### Congratulations, you’ve downloaded PLEX and run a docking tool!

